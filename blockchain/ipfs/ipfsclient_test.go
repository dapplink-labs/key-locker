package ipfs

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepo(t *testing.T) {
	str, err := CreateRepo()
	assert.NoError(t, err)
	fmt.Println(str)
}

func TestClient_AddFile(t *testing.T) {
	ctx, repoPath := context.Background(), ""
	repoPath, err := CreateRepo()
	assert.NoError(t, err)
	fmt.Printf("repoPath is %s", repoPath)

	c, err := New(ctx, []string{
		"/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWLXzGF1pXMYsNgv7yCGKdRwPqF2WAuqJwNFWvP2h3fNSp",
	}, repoPath)
	assert.NoError(t, err)

	cid, err := c.AddFile(ctx, []byte("test for add file"))
	assert.NoError(t, err)

	fmt.Println(cid)
}

func TestClient_GetFile(t *testing.T) {
	ctx, repoPath := context.Background(), ""
	repoPath, err := CreateRepo()
	assert.NoError(t, err)
	fmt.Printf("repoPath is %s", repoPath)

	c, err := New(ctx, []string{
		"/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWLXzGF1pXMYsNgv7yCGKdRwPqF2WAuqJwNFWvP2h3fNSp",
	}, repoPath)
	assert.NoError(t, err)

	// from local
	ret, err := c.GetFile(ctx, "QmbwC2qx2EXWeujvXZUKPWRyjfaUVA9s2ivTW28UQ9EZqM")
	assert.NoError(t, err)
	fmt.Println(string(ret))

	// from network
	ret, err = c.GetFile(ctx, "Qma71JMRwZc2aVMZ5McmbggfTgMJJQ8k3HKM8GpMeBR2CU")
	assert.NoError(t, err)
	fmt.Println(string(ret))
}
