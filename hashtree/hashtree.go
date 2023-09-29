package hashtree

import "crypto/sha256"

type Pair struct {
	mLHash [32]byte
	mRHash [32]byte
}

type Node struct {
	mHash      [32]byte
	mDirection bool
}

type Tree struct {
	mTreeMatrix [][][32]byte
}

// GenerateTree - will generate a MerkleTree, and store the tree in the t variable that is passed as a pointer.
// It will return the new hash of the block.
func GenerateTree(l [][32]byte, t *Tree) [32]byte {
	var newList [][32]byte

	if len(l)%2 == 1 {
		l = append(l, l[len(l)-1])
	}

	for i := 0; i < len(l); i += 2 {
		p := Pair{mLHash: l[i], mRHash: l[i+1]}

		newList = append(newList, p.generateHash())
	}

	if len(newList) == 1 {
		t.mTreeMatrix = append(t.mTreeMatrix, newList)
		return newList[0]
	} else if len(newList)%2 == 1 {
		newList = append(newList, newList[len(newList)-1])
	}
	t.mTreeMatrix = append(t.mTreeMatrix, newList)

	return GenerateTree(newList, t)
}

func ValidateData(name string, t *Tree, rootHash [32]byte) bool {
	nameHash := sha256.Sum256([]byte(name))

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

	return nameHash == rootHash
}

func (t *Tree) ClearTree() {
	t.mTreeMatrix = nil
}
