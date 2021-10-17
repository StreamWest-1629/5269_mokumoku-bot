package bot

import (
	"errors"
	"math/rand"
	"strconv"
)

type (
	GroupConn interface {
		GetWholeChats() *WholeChats
		MakeTextChat(name, topic string) (TextConn, error)
		MakeVoiceChat(name string) (VoiceConn, error)
	}

	Branches []Branch

	Branch struct {
		TextConn
		VoiceConn
	}
)

const MaxBranchMembers = 5

func SpreadBranches(conn GroupConn) (branches Branches, err error) {

	const BranchTopic = "やすみじかんは17分です。人とのんびりはなしながらやすみましょう。作業時間になったら自動で戻ります。会話記録は残りません。"

	// get mokumoku members
	whole := conn.GetWholeChats()
	memberIds := whole.MokuMoku.JoinMemberIds()

	// shuffle
	rand.Shuffle(len(memberIds), func(i, j int) { memberIds[i], memberIds[j] = memberIds[j], memberIds[i] })

	// make branches
	branches = make(Branches, (len(memberIds)-1)/MaxBranchMembers+1)
	for i := range branches {
		name := "やすみ_" + strconv.Itoa(i+1)

		// make branch
		if text, err := conn.MakeTextChat(name, BranchTopic); err != nil {
			return nil, errors.New("cannot make text chat: " + err.Error())
		} else if voice, err := conn.MakeVoiceChat(name); err != nil {
			return nil, errors.New("cannot make voice chat: " + err.Error())
		} else {
			branches[i].TextConn, branches[i].VoiceConn = text, voice
		}

		// make private
		if err := branches[i].MakePrivate(); err != nil {
			return nil, errors.New("cannot make chat private: " + err.Error())
		}
	}

	// move member to branch chats
	for i := range memberIds {
		if err := branches[i%len(branches)].MakeAllowance(memberIds[i]); err != nil {
			return nil, err
		} else if err := branches[i%len(branches)].MoveToHere(memberIds[i]); err != nil {
			return nil, errors.New("cannot move member to voice chat: " + err.Error())
		}
	}

	return branches, nil
}

func (b Branches) ClearBranches(conn GroupConn) error {

	// get mokumoku branch
	whole := conn.GetWholeChats()
	mokumoku := whole.MokuMoku

	// move to mokumoku room
	for i := range b {
		members := b[i].JoinMemberIds()
		for j := range members {
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
