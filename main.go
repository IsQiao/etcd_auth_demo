package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/auth/authpb"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"log"
	"time"
)

func main() {
	rootPwd := "123456"
	guestPwd := "123456"

	const (
		rootUsername  = "root"
		guestUsername = "guest"
		guestRole     = "guest_role"
		rootRole      = "root"
	)

	dbClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
		Username:    rootUsername,
		Password:    rootPwd,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.TODO()

	// check user existed
	getUserResp, err := dbClient.UserGet(ctx, "balabalalaaa")
	if err != nil {
		if err == rpctypes.ErrUserNotFound {
			fmt.Println("user not exist")
		} else {
			log.Fatalln(err)
		}
	}
	fmt.Println(getUserResp)

	if rootPwd != "" {
		getUserResp, err := dbClient.UserGet(ctx, rootUsername)
		if err != nil {
			if err == rpctypes.ErrUserNotFound {
				addUserResp, err := dbClient.UserAdd(ctx, rootUsername, rootPwd)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(addUserResp.Header)

				rootGrantRoleResp, err := dbClient.UserGrantRole(ctx, rootUsername, rootRole)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(rootGrantRoleResp.Header)
			} else {
				log.Fatalln(err)
			}
		}
		fmt.Println(getUserResp)
	}

	if guestPwd != "" {
		roleGetResp, err := dbClient.RoleGet(ctx, guestRole)
		if err != nil {
			if err == rpctypes.ErrRoleNotFound {
				roleAddResp, err := dbClient.RoleAdd(ctx, guestRole)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(roleAddResp)
			}
		}
		fmt.Println(roleGetResp)

		guestGetResp, err := dbClient.UserGet(ctx, guestUsername)
		if err != nil {
			if err == rpctypes.ErrUserNotFound {
				guestAddResp, err := dbClient.UserAdd(ctx, guestUsername, guestPwd)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(guestAddResp)
			}

			guestGrantRoleResp, err := dbClient.UserGrantRole(ctx, guestUsername, guestRole)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(guestGrantRoleResp.Header)

			dbClient.RoleGrantPermission(ctx, guestRole, "/", "/sys", clientv3.PermissionType(authpb.READWRITE))
			dbClient.RoleGrantPermission(ctx, guestRole, "/sys~", "0", clientv3.PermissionType(authpb.READWRITE))
		}

		fmt.Println(guestGetResp)

		authResp, err := dbClient.AuthEnable(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(authResp)
	}
}
