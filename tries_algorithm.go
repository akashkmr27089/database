package main

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	value    string
}

type Trie struct {
	root *TrieNode
}

type KeyValue struct {
	Key   string
	Value string
}

func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

func (t *Trie) Insert(word string, value string) {
	node := t.root
	for _, char := range word {
		if node.children[char] == nil {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.isEnd = true
	node.value = value
}

func (t *Trie) Search(word string) bool {
	node := t.root
	for _, char := range word {
		if node.children[char] == nil {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	node := t.root
	for _, char := range prefix {
		if node.children[char] == nil {
			return false
		}
		node = node.children[char]
	}
	return true
}

func (t *Trie) SearchPartial(prefix string) []KeyValue {
	node := t.root
	result := []KeyValue{}
	for _, char := range prefix {
		if node.children[char] == nil {
			return result
		}
		node = node.children[char]
	}
	t.searchWords(node, prefix, &result)
	return result
}

func (t *Trie) searchWords(node *TrieNode, prefix string, result *[]KeyValue) {
	if node.isEnd {
		final := KeyValue{Key: prefix, Value: node.value}
		*result = append(*result, final)
	}
	for char, child := range node.children {
		t.searchWords(child, prefix+string(char), result)
	}
}

func (t *Trie) Update(key, newValue string) {
	t.Delete(key)
	t.Insert(key, newValue)
}

func (t *Trie) Delete(word string) {
	t.deleteHelper(t.root, word, 0)
}

func (t *Trie) deleteHelper(node *TrieNode, word string, index int) bool {
	if index == len(word) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		return len(node.children) == 0
	}
	char := rune(word[index])
	child, ok := node.children[char]
	if !ok {
		return false
	}
	shouldDelete := t.deleteHelper(child, word, index+1)
	if shouldDelete {
		delete(node.children, char)
		return len(node.children) == 0 && !node.isEnd
	}
	return false
}

// func main2() {
// 	trie := NewTrie()
// 	words := []string{"apple", "banana", "orange", "app", "apps"}
// 	for _, word := range words {
// 		strconv.Itoa(123)
// 		trie.Insert(word, "221"+strconv.Itoa(rand.Int()))
// 	}

// 	fmt.Println(trie.Search("apple"))      // true
// 	fmt.Println(trie.Search("app"))        // true
// 	fmt.Println(trie.Search("oranges"))    // false
// 	fmt.Println(trie.StartsWith("or"))     // true
// 	fmt.Println(trie.StartsWith("xyz"))    // false
// 	fmt.Println(trie.SearchPartial("app")) // [app apple apps]
// 	fmt.Println(trie.SearchPartial("b"))   // [banana]
// 	fmt.Println(trie.SearchPartial("or"))  // [orange]
// 	fmt.Println(trie.SearchPartial("xyz")) // []

// 	trie.Update("apple", "mango")
// 	fmt.Println(trie.SearchPartial("apple"))
// 	fmt.Println(trie.Search("apple")) // false
// 	fmt.Println(trie.Search("mango")) // true

// 	trie.Delete("banana")
// 	fmt.Println(trie.Search("banana")) // false
// }
