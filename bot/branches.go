package bot

import "errors"

type (
	Rooms struct {
		Members
		MokuMoku VoiceChatConn
		branches []Branch
	}

	Branch struct {
		TextChatConn
		VoiceChatConn
		MokuMokuRoomId
		Name string
	}
)

const (
	MaxBreakingRoomMember = 4
)

func (r *Rooms) MakeBranch(name string, conn GroupConn) (MokuMokuRoomId, error) {

	// make id
	id := MokuMokuRoomId(len(r.branches) + 1)

	// make branch
	if branch, err := conn.MakeBranch(name); err != nil {
		return 0, errors.New("cannot make branch channel: " + err.Error())
	} else {
		// set name
		branch.Name, branch.MokuMokuRoomId = name, id

		// register
		r.branches = append(r.branches, branch)
		return id, nil
	}
}

func (r *Rooms) SpreadBranch() error {

	members := []Member{}

	// search move members
	for key := range r.Members {
		if member := r.Members[key]; member.Exists == RoomIDMokuMoku {
			members = append(members, member)
		}
	}

	if len(members) > 0 {

		// make branch channels
		for i, l := 0, (len(members) -1)/
		// TODO: WRITE HERE
		
		for i := range members {
			member[i].Allows = 
		}
	} else {
		return errors.New("cannot continue due to no one here")
	}

}

func (r *Rooms) ClearBranch() {
	for _, member := range r.Members {
		if member.Exists > RoomIDDisconnect {
			// move to mokumoku
			member.Allows = RoomIDMokuMoku
			r.MokuMoku.MoveToHere(member)
		} else {
			// erase member
			r.Members.RemoveMember(member)
		}
	}
	r.branches = r.branches[:0]
}
