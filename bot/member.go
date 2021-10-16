package bot

type (
	Members map[string]Member
	Member  struct {
		MemberConn
		Exists, Allows MokuMokuRoomId
	}
)

func (m Members) FindMember(id string) (member Member, exist bool) {
	member, exist = map[string]Member(m)[id]
	return
}

func (m Members) UpdateMember(member Member) {
	map[string]Member(m)[member.GetID()] = member
}

func (m Members) Joining(id string) (result bool) {
	_, result = map[string]Member(m)[id]
	return
}

func (m Members) RemoveMember(member Member) {
	delete(map[string]Member(m), member.GetID())
}
