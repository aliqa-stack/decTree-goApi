
package main

import (
    _"fmt"
    "math"
    "sort"
)

type Node struct {
    Leaf         bool
    PredictClass int
    Features     int
    Threshold    float64

    Left  *Node
    Right *Node
}

type Split struct {
    BestThreshold float64
    Feature       int
    Gain          float64
    Found         bool
}

func majority(y []int) int {
    count := make(map[int]int)

    for _, class := range y {
        count[class]++
    }

    bestCount := -1
    bestClass := 0

    for k, v := range count {
        if v > bestCount {
            bestCount = v
            bestClass = k
        }
    }

    return bestClass
}

func entropy(y []int) float64 {
    if len(y) == 0 {
        return 0
    }

    counts := make(map[int]int)

    for _, class := range y {
        counts[class]++
    }

    h := 0.0

    for _, value := range counts {
        p := float64(value) / float64(len(y))
        h -= p * math.Log2(p)
    }

    return h
}

func threshold(x []float64) []float64 {
    values := make([]float64, len(x))
    copy(values, x)

    sort.Float64s(values)

    result := []float64{}

    for i := 0; i < len(values)-1; i++ {
        if values[i] == values[i+1] {
            continue
        }

        mid := (values[i] + values[i+1]) / 2
        result = append(result, mid)
    }

    return result
}

func bestGain(y []int, x []float64) Split {
    best := Split{Gain: -1}

    if len(x) != len(y) || len(y) == 0 {
        return best
    }

    parent := entropy(y)

    for _, t := range threshold(x) {
        var leftY, rightY []int

        for i := range x {
            if x[i] <= t {
                leftY = append(leftY, y[i])
            } else {
                rightY = append(rightY, y[i])
            }
        }

        if len(leftY) == 0 || len(rightY) == 0 {
            continue
        }

        weighted :=
            float64(len(leftY))/float64(len(y))*entropy(leftY) +
            float64(len(rightY))/float64(len(y))*entropy(rightY)

        gain := parent - weighted

        if gain > best.Gain {
            best = Split{
                Found:         true,
                Gain:          gain,
                Feature:       0,
                BestThreshold: t,
            }
        }
    }

    return best
}

func TREE(y []int, x []float64, depth, maxDepth int) *Node {
    if len(y) == 0 || len(x) != len(y) {
        return nil
    }

    makeLeaf := func() *Node {
        return &Node{
            Leaf:         true,
            PredictClass: majority(y),
        }
    }

    if depth >= maxDepth {
        return makeLeaf()
    }

    pure := true

    for _, class := range y {
        if class != y[0] {
            pure = false
            break
        }
    }

    if pure {
        return makeLeaf()
    }

    split := bestGain(y, x)

    if !split.Found || split.Gain <= 0 {
        return makeLeaf()
    }

    var leftX, rightX []float64
    var leftY, rightY []int

    for i := range x {
        if x[i] <= split.BestThreshold {
            leftX = append(leftX, x[i])
            leftY = append(leftY, y[i])
        } else {
            rightX = append(rightX, x[i])
            rightY = append(rightY, y[i])
        }
    }

    node := &Node{
        Leaf:         false,
        PredictClass: majority(y),
        Features:     split.Feature,
        Threshold:    split.BestThreshold,
    }

    node.Left = TREE(leftY, leftX, depth+1, maxDepth)
    node.Right = TREE(rightY, rightX, depth+1, maxDepth)

    return node
}

func Predict(node *Node, x float64) int {
    if node.Leaf {
        return node.PredictClass
    }

    if x <= node.Threshold {
        return Predict(node.Left, x)
    }

    return Predict(node.Right, x)
}


