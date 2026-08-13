# Optimizing b & k To User Preferences


# Data & Optimization

## Data Format
so for each file we'll just have queries that should bring them up to the top. the closer to the top of the list they are, the better.
Meta data can be a toml file or whatever.
```toml
[filename]
queries = ["this is a query",
           "more queries"]
```

## Optimization Function
cost = sum of the position in the rank list of each element for each test query. position is normalized by length of document. start with linear growth in error and you could also try quadratic or smth like that (so that being at position 10 is much worse than position 5).

## Optimization Algorithm
Gradient free methods would be nice bc I don't really wanna do math rn. Either Baysian optimization or Simulated Annealing. Simulated Annealing would be easier I think.

## Optimization Pipeline:
what does your optimization pipeline look like?
so the optimization pipeline, the dirty gross way to do it:
for each query-doc pair in the training set (queries.toml) we run the query and get the position in overall returned documents of the desired document. that position is then normalized by the number of docs in the corpus. we want to minimize the total score from all of these ranking checks.


