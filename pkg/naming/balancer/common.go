package balancer

import (
	"strconv"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)

func getWeight(addr resolver.Address) int {
	if addr.Metadata == nil {
		return 1
	}
	md, ok := addr.Metadata.(*metadata.MD)
	if ok {
		values := md.Get("weight")
		if len(values) > 0 {
			w, err := strconv.Atoi(values[0])
			if err == nil {
				return w
			}
		}
	}
	return 1
}
