package hashtree

import "crypto/sha256"

// generateHash - generates the hash of a pair of nodes in a MerkleTree.
func (p *Pair) generateHash() [32]byte {
	var allBytes []byte
	allBytes = append(allBytes, p.mLHash[:]...)
	allBytes = append(allBytes, p.mRHash[:]...)
	return sha256.Sum256(allBytes)
}

// getListOfHashesToValidate - will return a slice with nodes inside the MerkleTree,
// which are used to compute the tree root
func getListOfHashesToValidate(index int, t *Tree) []Node {
	var newList []Node

	level := 0

	for {
		if index%2 == 0 {
			newList = append(newList, Node{mHash: t.mTreeMatrix[level][index+1], mDirection: true})
		} else {
			newList = append(newList, Node{mHash: t.mTreeMatrix[level][index-1], mDirection: false})
		}
		level++

		if len(t.mTreeMatrix[level]) == 1 {
			break
		}

		index = index / 2
	}

	return newList
}
