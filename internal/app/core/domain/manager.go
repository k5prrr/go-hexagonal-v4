package domain

import "math/big"

type UserGroupManager struct {
	User
	parent_id     big.Int
	shop_group_id int
}
