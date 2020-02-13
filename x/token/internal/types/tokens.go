package types

import (
	"encoding/json"
	"sort"
	"strings"

	linktype "github.com/line/link/types"
)

const (
	TokenTypeLength  = 4
	SmallestAlphanum = "0"
	LargestAlphanum  = "z"
	TokenIDLength    = linktype.TokenIDLen
	FungibleFlag     = SmallestAlphanum
	ReservedEmpty    = "0000"
	SmallestFTType   = "0001"
	SmallestNFTType  = "1001"
)

type Tokens []Token

func (ts Tokens) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (ts Tokens) IDAtIndex(index int) string { return ts[index].GetTokenID() }
func (ts Tokens) Len() int                   { return len(ts) }
func (ts Tokens) Less(i, j int) bool {
	return strings.Compare(ts[i].GetTokenID(), ts[j].GetTokenID()) == -1
}
func (ts Tokens) Swap(i, j int) { ts[i], ts[j] = ts[j], ts[i] }
func (ts Tokens) Sort() Tokens  { sort.Sort(ts); return ts }

func (ts Tokens) Append(tsB ...Token) Tokens {
	return append(ts, tsB...).Sort()
}

func (ts Tokens) Find(tokenID string) (Token, bool) {
	index := ts.find(tokenID)
	if index == -1 {
		return nil, false
	}
	return ts[index], true
}

func (ts Tokens) Update(token Token) (Tokens, bool) {
	index := ts.find(token.GetTokenID())
	if index == -1 {
		return ts, false
	}
	return append(append(ts[:index], token), ts[index+1:]...), true
}
func (ts Tokens) Remove(tokenID string) (Tokens, bool) {
	index := ts.find(tokenID)
	if index == -1 {
		return ts, false
	}
	return append(ts[:index], ts[index+1:]...), true
}

func (ts Tokens) Empty() bool {
	return ts.Len() == 0
}

func (ts Tokens) GetLatest() Token {
	if ts.Len() == 0 {
		return nil
	}

	return ts[ts.Len()-1]
}

func (ts Tokens) NextTokenID(prefix string) string {
	if len(prefix) > TokenIDLength {
		return ""
	}

	var tokens = ts
	if prefix != "" {
		tokens = ts.GetTokens(prefix)
	}
	latestToken := tokens.GetLatest()
	if latestToken == nil {
		return prefix + "0001"
		//return prefix + strings.Repeat(SmallestAlphanum, TokenIDLength-len(prefix))
	}
	nextTokenID := NextID(latestToken.GetTokenID(), prefix)
	for _, token := range tokens {
		if nextTokenID != token.GetTokenID() {
			return nextTokenID
		}
		nextTokenID = NextID(nextTokenID, prefix)
	}
	return ""
}

func (ts Tokens) NextTokenTypeForNFT() string {
	latestToken := ts.GetNFTs().GetLatest()
	if latestToken == nil {
		return SmallestNFTType
	}
	prefix := latestToken.GetTokenID()[:TokenTypeLength]
	for nextBaseID := NextID(prefix, ""); nextBaseID != prefix; nextBaseID = NextID(nextBaseID, "") {
		if nextBaseID[0] == FungibleFlag[0] {
			nextBaseID = "1" + nextBaseID[1:]
		}
		occupied := false
		ts.Iterate(nextBaseID, func(Token) bool { occupied = true; return true })
		if !occupied {
			return nextBaseID
		}
	}
	return ""
}

func (ts Tokens) NextTokenTypeForFT() string {
	latestToken := ts.GetFTs().GetLatest()
	if latestToken == nil {
		return SmallestFTType
	}

	prefix := latestToken.GetTokenID()[:TokenTypeLength]
	for nextBaseID := NextID(prefix, FungibleFlag); nextBaseID != prefix; nextBaseID = NextID(nextBaseID, FungibleFlag) {
		occupied := false
		ts.Iterate(nextBaseID, func(Token) bool { occupied = true; return true })
		if !occupied {
			return nextBaseID
		}
	}
	return ""
}

func (ts Tokens) GetTokens(prefix string) (tokens Tokens) {
	ts.Iterate(prefix, func(t Token) bool {
		tokens = append(tokens, t)
		return false
	})
	return tokens
}

func (ts Tokens) GetFTs() (tokens Tokens) {
	return ts.GetTokens(FungibleFlag)
}

func (ts Tokens) GetNFTs() (tokens Tokens) {
	ts.Iterate("", func(t Token) bool {
		if t.GetTokenID()[0] != FungibleFlag[0] {
			if t.GetTokenID()[TokenTypeLength:] != "0000" {
				tokens = append(tokens, t)
			}
		}
		return false
	})
	return tokens
}

func (ts Tokens) Iterate(prefix string, process func(Token) (stop bool)) {
	postLen := linktype.TokenIDLen - len(prefix)
	if postLen < 0 {
		return
	}

	start := prefix + strings.Repeat(SmallestAlphanum, postLen)
	end := prefix + strings.Repeat(LargestAlphanum, postLen)
	_, startIndex := BinarySearch(ts, start)
	if startIndex != -1 && strings.Compare(ts.IDAtIndex(startIndex), start) < 0 {
		startIndex++
	}
	_, endIndex := BinarySearch(ts, end)
	if endIndex != -1 && strings.Compare(ts.IDAtIndex(endIndex), end) > 0 {
		endIndex--
	}

	for index := startIndex; index >= 0 && index <= endIndex; index++ {
		if process(ts[index]) {
			return
		}
	}
}

func (ts Tokens) find(tokenID string) int {
	index, _ := BinarySearch(ts, tokenID)
	return index
}

func BinarySearch(group Findable, el string) (int, int) {
	if group.Len() == 0 {
		return -1, -1
	}
	low := 0
	high := group.Len() - 1
	median := 0
	for low <= high {
		median = (low + high) / 2
		switch compare := strings.Compare(group.IDAtIndex(median), el); {
		case compare == 0:
			// if group[median].element == el
			return median, median
		case compare == -1:
			// if group[median].element < el
			low = median + 1
		default:
			// if group[median].element > el
			high = median - 1
		}
	}
	return -1, median
}

func NextID(id string, prefix string) (nextTokenID string) {
	if len(prefix) >= len(id) {
		return prefix[:len(id)]
	}
	var toCharStr = "0123456789abcdefghijklmnopqrstuvwxyz"
	const toCharStrLength = 36 //int32(len(toCharStr))

	tokenIDInt := make([]int32, len(id))
	for idx, char := range id {
		if char >= '0' && char <= '9' {
			tokenIDInt[idx] = char - '0'
		} else {
			tokenIDInt[idx] = char - 'a' + 10
		}
	}
	for idx := len(tokenIDInt) - 1; idx >= 0; idx-- {
		char := tokenIDInt[idx] + 1
		if char < (int32)(toCharStrLength) {
			tokenIDInt[idx] = char
			break
		}
		tokenIDInt[idx] = 0
	}

	for _, char := range tokenIDInt {
		nextTokenID += string(toCharStr[char])
	}
	nextTokenID = prefix + nextTokenID[len(prefix):]

	return nextTokenID
}
