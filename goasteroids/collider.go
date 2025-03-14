package goasteroids

import "github.com/solarlune/resolv"

// checkCollision will check object obj for collisions. If against is nil,
// then we check obj against ALL other objects; if it is not, we only
// check against one object (against).
func (g *GameScene) checkCollision(obj, against *resolv.Circle) bool {
	if against == nil {
		return obj.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: obj.SelectTouchingCells(1).FilterShapes(),
			OnIntersect: func(set resolv.IntersectionSet) bool {
				return true
			},
		})
	}

	return obj.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: against.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			return true
		},
	})
}
