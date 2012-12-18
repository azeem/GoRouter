Router
======

Simple Generic request Router for Go Language

Example
-------

    routes := []*Route {
        NewRoute().Path("pages", "archive").Handle("First Router"),
        NewRoute().Path("abc", "def", "ghi").Handle("Second Router"),
    }

    if result := MatchRoute(routes, req); result != nil {
        fmt.Println("Matching handle = ", result.Handle)
    } else {
        fmt.Println("Match not found")
    }
