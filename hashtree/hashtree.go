package hashtree

import (
	"crypto/sha256"
	"github.com/TheJ0lly/GoChain/metadata"
)

type Pair struct {
	mLHash [32]byte
	mRHash [32]byte
}

type Node struct {
	mHash      [32]byte
	mDirection bool //if direction is true it means the hash is on the right, else is on the left
}

type Tree struct {
	mTreeMatrix [][][32]byte
	RootHash    [32]byte
}

// GenerateTree - will generate a MerkleTree, and store the tree in the t variable that is passed as a pointer.
// It will return the new hash of the block.
func GenerateTree(l [][32]byte) *Tree {
	if len(l)%2 == 1 {
		l = append(l, l[len(l)-1])
	}

	t := &Tree{
		mTreeMatrix: nil,
		RootHash:    [32]byte{},
	}

	t.mTreeMatrix = append(t.mTreeMatrix, l)

	generateTreeRecursive(l, t)

	return t
}

// generateTreeRecursive - will generate the tree matrix recursively
func generateTreeRecursive(l [][32]byte, t *Tree) {
	var newList [][32]byte

	if len(l)%2 == 1 {
		l = append(l, l[len(l)-1])
	}

	for i := 0; i < len(l); i += 2 {
		p := Pair{mLHash: l[i], mRHash: l[i+1]}

		newList = append(newList, p.generateHash())
	}

	t.mTreeMatrix = append(t.mTreeMatrix, newList)

	if len(newList) == 1 {
		t.RootHash = newList[0]
		return //Not sure if needed
	} else {
		generateTreeRecursive(newList, t)
	}

}

func ValidateData(name metadata.MetaData, t *Tree) bool {
	nameHash := sha256.Sum256([]byte(name.GetMetaDataString()))

	var l []Node

	for i, k := range t.mTreeMatrix[0] {
		if k == nameHash {
			l = getListOfHashesToValidate(i, t)
		}
	}

	var p *Pair

	for _, k := range l {
		if k.mDirection {
			p = &Pair{mLHash: nameHash, mRHash: k.mHash}
		} else {
			p = &Pair{mLHash: k.mHash, mRHash: nameHash}
		}
		nameHash = p.generateHash()
	}

	return nameHash == t.RootHash
}

// ClearTree - will clear the current tree
func (t *Tree) ClearTree() {
	t.mTreeMatrix = nil
}
