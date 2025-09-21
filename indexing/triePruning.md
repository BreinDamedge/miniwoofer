# Re: how tf do you prune a branch in this god forsaken data structure?

## when to prune:
when your node has nil docIds then it may be time to prune! if it also has nil or empty (len == 0 ) children then it is a leaf node so you can definitely prune it. The catch is that you can't just cut off the leaf, you have to remove the entire dead branch. So how do we do that?

## How to Prune:
using trajectories would be make life easier:
if your trajectory is a slice of nodes that could work
given a slice of node pointers (which is your path in the tree)
you can traverse them backward and prune as you go 
if children is empty then delete this node (by setting the parent's pointer to nil/removing it from the dict).

given a trajectory, walk backward and if a node is empty and has no children, delete it from the parent.

for a node in the trajectory you go if it's empty and should be deleted:
  leaf rune (r) and then continue 
