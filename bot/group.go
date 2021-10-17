package bot

import (
	"errors"
	"fmt"
	"strconv"
)

type (
	ChatID int

	WholeChats struct {
		Random   TextConn
		Todo     TextConn
		MokuMoku VoiceConn
	}

	GroupConn interface {
		Update() error
		MakeTextChat(name, description string) (TextConn, error)
		MakeVoiceChat(name string) (VoiceConn, error)
		RemoveChat(name string)
		Initialize() (whole *WholeChats, err error)
		GetMemberIds() ([]string, error)
		GetMember(memberId string) (MemberConn, error)
		VoiceState(memberId string) (chatId *string)
	}

	Member struct {
		MemberConn
		Exists, Allows ChatID
	}

	BranchGroup struct {
		GroupConn
		*WholeChats
		members  map[string]Member
		Branches []__Branch
		vcs      map[string]int
	}

	__Branch struct {
		TextConn
		VoiceConn
	}
)

const (
	RoomIDDisconnect ChatID = -iota
	RoomIDMokuMoku
	RoomIDOther
)

const (
	MaxBranchMembers = 4
)

func (g *BranchGroup) CheckMember(memberId string, joinChatId string) error {

	if member, exist := g.members[memberId]; exist {
		if idx, exist := g.vcs[joinChatId]; exist {
			if idx != int(member.Allows) {
				g.Branches[idx].VoiceConn.MoveToHere(memberId)
			}
		} else if joinChatId == g.MokuMoku.GetID() && member.Allows > 0 {
			g.Branches[idx].VoiceConn.MoveToHere(memberId)
		}
	} else {
		if _, exist := g.vcs[joinChatId]; exist {
			g.MokuMoku.MoveToHere(memberId)
		}
	}
	return nil
}

func (g *BranchGroup) ClearBranch() error {

	// check voice state
	for i := range g.members {
		if chatId := g.VoiceState(g.members[i].GetID()); chatId != nil {
			if _, exist := g.vcs[*chatId]; exist {
				g.MokuMoku.MoveToHere(g.members[i].GetID())
			}
		}
	}

	// remove chats
	for i := range g.Branches {
		g.RemoveChat(g.Branches[i].TextConn.GetID())
		g.RemoveChat(g.Branches[i].VoiceConn.GetID())
	}

	return nil

}

func SpreadBranch(conn GroupConn) (branches *BranchGroup, err error) {

	fmt.Println("begin spread branch")

	if whole, err := conn.Initialize(); err != nil {
		return nil, errors.New("cannot make group initialized: " + err.Error())
	} else {

		branches = &BranchGroup{
			GroupConn:  conn,
			WholeChats: whole,
			members:    map[string]Member{},
			Branches:   []__Branch{},
			vcs:        map[string]int{},
		}
	}

	memberIds, err := conn.GetMemberIds()
	if err != nil {
		return nil, err
	}

	// voice state
	for i := range memberIds {
		if chatId := branches.VoiceState(memberIds[i]); chatId != nil {
			if *chatId == branches.WholeChats.MokuMoku.GetID() {
				if mem, err := branches.GetMember(*chatId); err != nil {
					return nil, errors.New("cannot get member: " + err.Error())
				} else {
					branches.members[*chatId] = Member{
						MemberConn: mem,
					}
				}
			}
		}
	}

	// make branch channel
	for i, l := 0, (len(branches.members)-1)/MaxBranchMembers+1; i < l; i++ {

		name := "Breaking Room " + strconv.Itoa(i+1)
		if text, err := branches.MakeTextChat(name, ""); err != nil {
			return nil, errors.New("cannot make breaking text chat: " + err.Error())
		} else if err := text.MakePrivate(); err != nil {
			return nil, errors.New("cannot make breaking text chat private: " + err.Error())
		} else if voice, err := branches.MakeVoiceChat(name); err != nil {
			return nil, errors.New("cannot make breaking voice chat: " + err.Error())
		} else if err := voice.MakePrivate(); err != nil {
			return nil, errors.New("cannot make breaking voice chat private: " + err.Error())
		} else {
			branches.Branches = append(branches.Branches, __Branch{
				TextConn:  text,
				VoiceConn: voice,
			})
		}
	}

	// move member
	i := 0
	for key := range branches.members {
		j := i % len(branches.Branches)

		// get allowance
		if err := branches.Branches[j].TextConn.MakeMemberAllow(key); err != nil {
			return nil, err
		} else if err := branches.Branches[j].VoiceConn.MakeMemberAllow(key); err != nil {
			return nil, err
		}

		// member instance move
		member := branches.members[key]
		member.Allows = ChatID(j + 1)
		branches.members[key] = member

		//
		if err := branches.Branches[j].MoveToHere(key); err != nil {
			return nil, err
		}
		i++
	}

	for i := range branches.Branches {
		branches.vcs[branches.Branches[i].VoiceConn.GetID()] = i + 1
	}

	return branches, nil
}
