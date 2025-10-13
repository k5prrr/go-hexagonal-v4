package domain

import "math/big"

type UserGroupSeller struct {
	User
	parent_id big.Int
	//user_id       big.Int
	shop_group_id int
	shop_id       big.Int
}
