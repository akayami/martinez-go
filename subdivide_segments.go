package martinez_go

import (
	"math"
)

// Assume definitions of SweepEvent, compareSegments, computeFields, and possibleIntersection exist

// Subdivide processes the event queue for the sweep line algorithm.
func SubdivideSegments(eventQueue TinyQueue, subject, clipping [][][]Point, sbbox, cbbox *[4]float64, operation int) []*SweepEvent {
	// sweepLine := splay_tree_go.NewSplayTree[int](nil) // NewSweepLineTree()
	// sweepLine := NewSweepLineTree()
	sweepLine := NewSplayTree(CompareSegments)
	var sortedEvents []*SweepEvent

	rightbound := math.Min(sbbox[2], cbbox[2])

	var prevNode, nextNode, beginNode *Node

	for eventQueue.Len() > 0 {
		event := eventQueue.Pop() // heap.Pop(&eventQueue).(*SweepEvent)
		sortedEvents = append(sortedEvents, event)

		if (operation == Intersection && event.Point.X > rightbound) ||
			(operation == Difference && event.Point.X > sbbox[2]) {
			break
		}

		if event.Left {
			node := sweepLine.Insert(NewSplayTreeNode(event))
			prevNode = node
			nextNode = node
			beginNode = sweepLine.MinNode(sweepLine.Root)

			if prevNode != beginNode {
				prevNode = sweepLine.Prev(prevNode)
			} else {
				prevNode = nil
			}

			nextNode = sweepLine.Next(nextNode)

			var prevEvent, prevprevEvent *SweepEvent
			if prevNode != nil {
				prevEvent = prevNode.Key
			}

			ComputeFields(event, prevEvent, operation)

			if nextNode != nil {
				if PossibleIntersection(event, nextNode.Key, eventQueue) == 2 {
					ComputeFields(event, prevEvent, operation)
					ComputeFields(nextNode.Key, event, operation)
				}
			}

			if prevNode != nil {
				if PossibleIntersection(prevNode.Key, event, eventQueue) == 2 {
					var prevprevNode = prevNode
					if prevNode != beginNode {
						prevprevNode = sweepLine.Prev(prevprevNode)
					} else {
						prevprevNode = nil
					}
					if prevprevNode != nil {
						prevprevEvent = prevprevNode.Key
					}
					ComputeFields(prevEvent, prevprevEvent, operation)
					ComputeFields(event, prevEvent, operation)
				}
			}
		} else {
			event = event.OtherEvent
			node := sweepLine.Find(event)
			if node != nil {
				prevNode = node
				nextNode = node
				if prevNode != beginNode {
					prevNode = sweepLine.Prev(prevNode)
				} else {
					prevNode = nil
				}
				nextNode = sweepLine.Next(nextNode)

				sweepLine.Remove(event)
				if nextNode != nil && prevNode != nil {
					PossibleIntersection(prevNode.Key, nextNode.Key, eventQueue)
				}
			}
		}
	}
	return sortedEvents
}
