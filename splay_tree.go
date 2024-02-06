package martinez_go

// Key Define basic types for Key and Value
type Key *SweepEvent // Assuming Key is of type int for simplicity; adjust as needed.
type Value any

// Comparator function type
type Comparator func(a, b *SweepEvent) int

// Visitor function type for traversing nodes
type Visitor func(node *Node)
type RangeVisitor func(node *Node) bool

// NodePrinter function type for printing nodes
type NodePrinter func(node *Node) string

// Node structure
type Node struct {
	Key   *SweepEvent
	Value any
	Left  *Node
	Right *Node
	Next  *Node
}

func NewSplayTreeNode(key *SweepEvent) *Node {
	return &Node{Key: key}
}

// SplayTree structure
type SplayTree struct {
	Root       *Node
	Size       int
	Comparator Comparator
}

// NewTree creates a new instance of SplayTree with a given comparator
func NewSplayTree(comparator Comparator) *SplayTree {
	if comparator == nil {
		comparator = CompareSegments // Ensure there's a default comparator
	}
	return &SplayTree{
		Comparator: comparator,
	}
}

// Insert Inserts a key, allows duplicates
func (tree *SplayTree) Insert(node *Node) *Node {
	tree.Root = insert(node.Key, nil, tree.Root, tree.Comparator)
	tree.Size++
	return tree.Root
}

// Add inserts a key and its associated data into the tree if the key does not already exist.
// It returns the new root of the tree after insertion.
func (tree *SplayTree) Add(key Key, value Value) *Node {
	node := &Node{
		Key:   key,
		Value: value,
		Left:  nil,
		Right: nil,
	}

	if tree.Root == nil {
		// If the tree is empty, simply insert the new node and return it.
		tree.Root = node
		tree.Size++
		return tree.Root
	}

	// Splay the tree around the key to bring the closest node to the root.
	t := splay(key, tree.Root, tree.Comparator)

	// Compare the new key with the root's key after splaying.
	cmp := tree.Comparator(key, t.Key)
	if cmp == 0 {
		// If the key is already in the tree, just set t as the new root.
		tree.Root = t
	} else {
		// Insert the new node in the appropriate position.
		if cmp < 0 {
			node.Left = t.Left
			node.Right = t
			t.Left = nil
		} else { // cmp > 0
			node.Right = t.Right
			node.Left = t
			t.Right = nil
		}
		tree.Size++
		tree.Root = node
	}

	return tree.Root
}

// Remove is a public method that removes a node with the given key from the tree, if it exists.
func (tree *SplayTree) Remove(key Key) {
	tree.Root = removeNode(key, tree.Root, tree.Comparator, &tree.Size)
}

// removeNode is a private function that deletes a node with the given key from the tree, if it exists.
func removeNode(i Key, t *Node, comparator Comparator, size *int) *Node {
	if t == nil {
		return nil
	}

	t = splay(i, t, comparator)
	cmp := comparator(i, t.Key)
	if cmp == 0 { // Found the node to remove
		var x *Node
		if t.Left == nil {
			x = t.Right
		} else {
			// Splay the left subtree to bring the maximum element to the top
			x = splay(i, t.Left, comparator)
			x.Right = t.Right
		}
		*size-- // Decrement the size of the tree
		return x
	}
	return t // Node wasn't found; return the tree as is
}

// Pop removes and returns the node with the smallest key.
func (tree *SplayTree) Pop() (Key, Value, bool) {
	if tree.Root == nil {
		// Return zero values and false to indicate that the tree is empty and no node was popped.
		return nil, nil, false
	}

	node := tree.Root
	for node.Left != nil {
		node = node.Left
	}

	// Splay the tree around the node with the smallest key to bring it to the root.
	tree.Root = splay(node.Key, tree.Root, tree.Comparator)
	// Remove the node with the smallest key.
	tree.Root = removeNode(node.Key, tree.Root, tree.Comparator, &tree.Size)

	// Return the key and data of the removed node, and true to indicate success.
	return node.Key, node.Value, true
}

// FindStatic searches for a node by key without splaying.
func (tree *SplayTree) FindStatic(key Key) (*Node, bool) {
	current := tree.Root
	for current != nil {
		cmp := tree.Comparator(key, current.Key)
		if cmp == 0 {
			// Node found, return the node and true to indicate success.
			return current, true
		} else if cmp < 0 {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	// Node not found, return nil and false.
	return nil, false
}

// Find searches for a node by key and splays the tree.
// If the node is found, it becomes the new root of the tree.
// If the node is not found, the last accessed node becomes the new root.
func (tree *SplayTree) Find(key Key) *Node {
	if tree.Root != nil {
		tree.Root = splay(key, tree.Root, tree.Comparator)
		if tree.Comparator(key, tree.Root.Key) != 0 {
			// Node with the given key was not found.
			return nil
		}
	}
	// Return the root (either the found node or the last accessed node)
	return tree.Root
}

// Contains checks if the tree contains a node with the specified key.
func (tree *SplayTree) Contains(key Key) bool {
	current := tree.Root
	for current != nil {
		cmp := tree.Comparator(key, current.Key)
		if cmp == 0 {
			// Key found in the tree.
			return true
		} else if cmp < 0 {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	// Key not found in the tree.
	return false
}

// Keys returns a slice of all keys in the tree.
func (tree *SplayTree) Keys() []Key {
	var keys []Key
	tree.ForEach(func(node *Node) {
		keys = append(keys, node.Key)
	})
	return keys
}

// Values returns a slice of all data in the tree nodes.
func (tree *SplayTree) Values() []Value {
	var values []Value
	tree.ForEach(func(node *Node) {
		values = append(values, node.Value)
	})
	return values
}

// Min returns the smallest key in the tree or zero value if the tree is empty.
func (tree *SplayTree) Min() (Key, bool) {
	minNode := tree.MinNode(tree.Root)
	if minNode != nil {
		return minNode.Key, true
	}
	return nil, false // Assuming Key is int for simplicity; adjust the zero value as necessary.
}

// Max returns the largest key in the tree or zero value if the tree is empty.
func (tree *SplayTree) Max() (Key, bool) {
	maxNode := tree.MaxNode(tree.Root)
	if maxNode != nil {
		return maxNode.Key, true
	}
	return nil, false // Adjust the zero value as necessary.
}

// ForEach traverses the tree in-order and applies the visitor function to each node.
func (tree *SplayTree) ForEach(visitor Visitor) {
	if tree.Root == nil {
		return
	}

	current := tree.Root
	var Q []*Node // Stack for in-order traversal
	done := false

	for !done {
		if current != nil {
			Q = append(Q, current) // Push current onto stack
			current = current.Left // Move to left child
		} else {
			if len(Q) > 0 {
				current = Q[len(Q)-1] // Pop from stack
				Q = Q[:len(Q)-1]

				visitor(current) // Visit the node

				current = current.Right // Move to right child
			} else {
				done = true
			}
		}
	}
}

// Range performs an in-order traversal of the tree within a key range from low to high.
// If the visitor function returns true, the traversal stops.
func (tree *SplayTree) Range(low, high Key, fn RangeVisitor) {
	var Q []*Node // Stack for in-order traversal
	node := tree.Root

	for len(Q) > 0 || node != nil {
		if node != nil {
			Q = append(Q, node) // Push current onto stack
			node = node.Left    // Move to left child
		} else {
			node = Q[len(Q)-1] // Pop from stack
			Q = Q[:len(Q)-1]

			// Compare node key with high; if node's key is greater, stop the traversal.
			if tree.Comparator(node.Key, high) > 0 {
				break
			}
			// If node's key is within the range, visit the node.
			if tree.Comparator(node.Key, low) >= 0 {
				if fn(node) { // If fn returns true, stop the traversal.
					return
				}
			}
			node = node.Right // Move to right child
		}
	}
}

// MinNode returns the node with the smallest key starting from the given node.
func (tree *SplayTree) MinNode(t *Node) *Node {
	if t == nil {
		return nil
	}
	for t.Left != nil {
		t = t.Left
	}
	return t
}

// MaxNode returns the node with the largest key starting from the given node.
func (tree *SplayTree) MaxNode(t *Node) *Node {
	if t == nil {
		return nil
	}
	for t.Right != nil {
		t = t.Right
	}
	return t
}

// At returns the node at the given index in an in-order traversal of the tree.
// If the index is out of bounds, it returns nil.
func (tree *SplayTree) At(index int) *Node {
	if tree.Root == nil {
		return nil // SplayTree is empty
	}

	current := tree.Root
	var Q []*Node // Stack for in-order traversal
	i := 0

	for {
		if current != nil {
			Q = append(Q, current) // Push current onto stack
			current = current.Left // Move to left child
		} else {
			if len(Q) == 0 {
				break // All nodes visited
			}
			current = Q[len(Q)-1] // Pop from stack
			Q = Q[:len(Q)-1]

			if i == index {
				return current // Found the node at the specified index
			}
			i++
			current = current.Right // Move to right child
		}
	}

	return nil // Index is out of bounds
}

// Next returns the successor of the given node in the tree.
func (tree *SplayTree) Next(d *Node) *Node {
	var successor *Node
	root := tree.Root

	if d.Right != nil {
		successor = d.Right
		for successor.Left != nil {
			successor = successor.Left
		}
		return successor
	}

	for root != nil {
		cmp := tree.Comparator(d.Key, root.Key)
		if cmp == 0 {
			break
		} else if cmp < 0 {
			successor = root
			root = root.Left
		} else {
			root = root.Right
		}
	}

	return successor
}

// Prev returns the predecessor of the given node in the tree.
func (tree *SplayTree) Prev(d *Node) *Node {
	var predecessor *Node
	root := tree.Root

	if d.Left != nil {
		predecessor = d.Left
		for predecessor.Right != nil {
			predecessor = predecessor.Right
		}
		return predecessor
	}

	for root != nil {
		cmp := tree.Comparator(d.Key, root.Key)
		if cmp == 0 {
			break
		} else if cmp > 0 {
			predecessor = root
			root = root.Right
		} else {
			root = root.Left
		}
	}

	return predecessor
}

// Clear removes all nodes from the tree.
func (tree *SplayTree) Clear() {
	tree.Root = nil
	tree.Size = 0
}

// toList helper function to convert tree to list, assuming definition is given.
func toList(node *Node) []*Node {
	// Assume toList implementation is defined here.
	return nil // Placeholder
}

// ToList converts the tree to a list.
func (tree *SplayTree) ToList() []*Node {
	return toList(tree.Root)
}

func merge(left, right *Node, comparator Comparator) *Node {
	if right == nil {
		return left
	}
	if left == nil {
		return right
	}

	// Splay the right tree around the maximum key in the left tree.
	right = splay(left.Key, right, comparator)

	// Attach the left tree as the left child of the right tree's root.
	right.Left = left

	return right
}

func insert(i Key, data Value, t *Node, comparator Comparator) *Node {
	node := &Node{
		Key:   i,
		Value: data,
		Left:  nil,
		Right: nil,
	}

	if t == nil {
		// If the tree is empty, return the new node as the root of the tree.
		return node
	}

	// Splay the tree around the key i.
	t = splay(i, t, comparator)

	// Compare the new node's key with the root's key.
	cmp := comparator(i, t.Key)
	if cmp < 0 {
		// If the new node's key is less than the root's key, insert the new node to the left of the root.
		node.Left = t.Left
		node.Right = t
		t.Left = nil
	} else if cmp >= 0 {
		// If the new node's key is greater than or equal to the root's key, insert the new node to the right of the root.
		node.Right = t.Right
		node.Left = t
		t.Right = nil
	}

	// Return the new node, which is now the root of the splayed tree.
	return node
}

func splay(i Key, t *Node, comparator Comparator) *Node {
	if t == nil {
		return nil
	}

	N := &Node{} // Temporary node for simplifying the tree manipulation
	l, r := N, N

	for {
		cmp := comparator(i, t.Key)
		if cmp < 0 {
			if t.Left == nil {
				break
			}
			if comparator(i, t.Left.Key) < 0 {
				y := t.Left // Rotate right
				t.Left = y.Right
				y.Right = t
				t = y
				if t.Left == nil {
					break
				}
			}
			r.Left = t // Link right
			r = t
			t = t.Left
		} else if cmp > 0 {
			if t.Right == nil {
				break
			}
			if comparator(i, t.Right.Key) > 0 {
				y := t.Right // Rotate left
				t.Right = y.Left
				y.Left = t
				t = y
				if t.Right == nil {
					break
				}
			}
			l.Right = t // Link left
			l = t
			t = t.Right
		} else {
			break
		}
	}

	// Assemble
	l.Right = t.Left
	r.Left = t.Right
	t.Left = N.Right
	t.Right = N.Left

	return t
}
