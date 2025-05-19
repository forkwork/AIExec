package real

func (r *RealBackwardsInvocation) aiexecPath(path ...string) string {
	path = append([]string{"inner", "api"}, path...)
	return r.aiexecInnerApiBaseurl.JoinPath(path...).String()
}
