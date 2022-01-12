package bot

import (
	"errors"
	"math/rand"
	"strconv"
)

type (
	GroupConn interface {
		MakePrivateTextChat(name, topic string, allowMemberIds []string) (TextConn, error)
		MakePrivateVoiceChat(name string, allowMemberIds []string) (VoiceConn, error)
		SetStateMessage(message string)
		MemberMute(memberId string, mute bool)
	}

	Branches []Branch

	Branch struct {
		TextConn
		VoiceConn
	}
)

const MaxBranchMembers = 5

func SpreadBranches(conn GroupConn, args *EventArgs) (branches Branches, err error) {

	const BranchTopic = "やすみじかんは17分です。人とのんびりはなしながらやすみましょう。作業時間になったら自動で戻ります。会話記録は残りません。"

	// get mokumoku members
	memberIds := []string{}

	// check member is ignored or not
	for i, ids := 0, args.MokuMoku.JoinMemberIds(); i < len(ids); i++ {
		if _, exist := args.BranchIgnore[ids[i]]; !exist {
			memberIds = append(memberIds, ids[i])
		}
	}

	// shuffle
	rand.Shuffle(len(memberIds), func(i, j int) { memberIds[i], memberIds[j] = memberIds[j], memberIds[i] })
	allowMemberIds := [][]string{}

	// make branches instance
	branches = make(Branches, (len(memberIds)-1)/MaxBranchMembers+1)

	for range branches {
		allowMemberIds = append(allowMemberIds, []string{})
	}

	for i := range memberIds {
		allowMemberIds[i%len(branches)] = append(allowMemberIds[i%len(branches)], memberIds[i])
	}

	for i := range branches {
		name := "やすみ_" + strconv.Itoa(i+1)

		// make branch
		if text, err := conn.MakePrivateTextChat(name, BranchTopic, allowMemberIds[i]); err != nil {
			return nil, errors.New("cannot make text chat: " + err.Error())
		} else if voice, err := conn.MakePrivateVoiceChat(name, allowMemberIds[i]); err != nil {
			return nil, errors.New("cannot make voice chat: " + err.Error())
		} else {
			branches[i].TextConn, branches[i].VoiceConn = text, voice
		}
		// move member to branch chats
		for j := range allowMemberIds[i] {
			if err := branches[j%len(branches)].MoveToHere(allowMemberIds[i][j]); err != nil {
				return nil, errors.New("cannot move member to voice chat: " + err.Error())
			}
		}
	}

	return branches, nil
}

func (b Branches) ClearBranches(conn GroupConn, args *EventArgs) error {

	// get mokumoku branch
	mokumoku := args.MokuMoku

	// move to mokumoku room
	for i := range b {
		members := b[i].JoinMemberIds()

		for j := range members {
			conn.MemberMute(members[i], true)
			if err := mokumoku.MoveToHere(members[j]); err != nil {
				return err
			}
		}
		b[i].TextConn.Delete()
		b[i].VoiceConn.Delete()
	}
	return nil
}

func (b *Branch) MakePrivate() error {
	if err := b.TextConn.MakePrivate(); err != nil {
		return errors.New("cannot make text chat private: " + err.Error())
	} else if err := b.VoiceConn.MakePrivate(); err != nil {
		return errors.New("cannot make voice chat private: " + err.Error())
	}
	return nil
}

func (b *Branch) MakeAllowance(memberId string) error {
	if err := b.TextConn.AllowAccess(memberId); err != nil {
		return errors.New("cannot make a member's text chat allowance: " + err.Error())
	} else if err := b.VoiceConn.AllowAccess(memberId); err != nil {
		return errors.New("cannot make a member's voice chat allowance: " + err.Error())
	}
	return nil
}
