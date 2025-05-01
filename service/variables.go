package service

import "errors"

var CommitID = "CommitID"
var NodeID = "NodeID"
var ErrOrgNotFound = errors.New("organization not found")
var ErrOwnerNotFound = errors.New("owner not found")
var ErrRepoNotFound = errors.New("repo not found")
var ErrBranchesNotFound = errors.New("branch not found")
var ErrRepoAlreadyExists = errors.New("repo name already exists")
var ErrBranchesAlreadyExists = errors.New("branch name already exists")
var ErrPRNotFound = errors.New("no PRs found. Invalid PR number")
var ErrOwnerNotInSameOrg = errors.New("owners are not belonging to same organization")
var ErrPRAlreadyExists = errors.New("already PR exists on the branch")
var ErrPRAlreadyClosed = errors.New("already PR got closed on the branch")
var ErrInvalidBranchName = errors.New("invalid branch name. Specify as refs/heads/<branch>")
