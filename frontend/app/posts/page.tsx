"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

import { useState, useEffect, useRef, useCallback } from "react";
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
  Container,
  Typography,
  Stack,
  Paper,
  Button,
  Select,
  MenuItem,
  FormControl,
  Box,
} from "@mui/material";
import { useRouter, useSearchParams } from "next/navigation";

const theme = createTheme({
  palette: {
    primary: { main: "#6366f1" },
    secondary: { main: "#8b5cf6" },
    background: { default: "#f8fafc", paper: "#ffffff" },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
  },
});

type Post = {
  id: number;
  title: string;
  description: string;
  likes: number;
  dislikes: number;
  views: number;
};

export default function PostsPage() {
  const router = useRouter();

  const [posts, setPosts] = useState<Post[]>([]);
  const [sortBy, setSortBy] = useState<"popularity" | "recency" | "views">("popularity");
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [isResetting, setIsResetting] = useState(false);

  const containerRef = useRef<HTMLDivElement>(null);
  const [query, setQuery] = useState<string | null>(null);
  const [rank, setRank] = useState<"popularity" | "recency" | "views" | null>(null);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    setQuery(params.get("q"));
    const sortParam = params.get("sort_by") as "popularity" | "recency" | "views" | null;
    if (sortParam) setRank(sortParam);
  }, []);

  const fetchPosts = async () => {
    if (!hasMore || loading) return;

    setLoading(true);
    try {
        const path = query ? `${apiUrl}/public/topics/0/posts/search?page=${page}&q=${query}&sort_by=${sortBy}` 
                           : `${apiUrl}/public/posts?page=${page}&sort_by=${sortBy}`;
        const res = await fetch(
        `${path}`
        );
        const json = await res.json();

        if (json.posts.length === 0) {
        setHasMore(false);
        return;
        }

        setPosts((prev) => [...prev, ...json.posts]);

        if (json.posts.length < 10) setHasMore(false);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  //pagination and sort
  useEffect(() => {
    setIsResetting(true);
    setPosts([]);
    setPage(1);
    setHasMore(true);
  }, [sortBy, query]);

    useEffect(() => {
    if (rank) setSortBy(rank);
    }, [rank]);

  useEffect(() => {
    if (isResetting) {
      fetchPosts();
      setIsResetting(false);
    } 
    else if(page !== 1) {
      fetchPosts();
    }
  }, [page, isResetting]);

  const handleScroll = useCallback(() => {
    if (!containerRef.current || loading || !hasMore) return;

    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    if (scrollHeight - scrollTop <= clientHeight + 100) {
      setPage((prev) => prev + 1);
    }
  }, [loading, hasMore]);

  const handleViewPost = async (postID: number, views: number, likes: number, dislikes: number) => {
    try {
      await fetch(`${apiUrl}/public/posts/${postID}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ 
          "views": views,
          "likes": likes,
          "dislikes": dislikes
        }),
        credentials: "include",
      });

      router.push(`/posts/${postID}`);
    } catch (err) {
      console.error("Error updating view:", err);
    }
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Container
        maxWidth="md"
        sx={{ mt: 4, height: "80vh", overflowY: "auto" }}
        ref={containerRef}
        onScroll={handleScroll}
      >
        {/* Header + Sort */}
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h4" sx={{ fontWeight: 700 }}>
            All Posts
          </Typography>

          {/* Sort dropdown */}
          <FormControl
            size="small"
            sx={{ minWidth: 120, ml: 2, background: "white", borderRadius: 1 }}
          >
            <Select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as "popularity" | "recency" | "views")}
              displayEmpty
              sx={{ "& .MuiSelect-select": { py: 0.5, px: 1.5 } }}
            >
              <MenuItem value="popularity"> Popularity</MenuItem>
              <MenuItem value="recency"> Recency</MenuItem>
              <MenuItem value="views"> Views</MenuItem>
            </Select>
          </FormControl>
        </Box>

        <Stack spacing={3}>
          {posts.length === 0 && loading ? (
            <Typography>Loading posts...</Typography>
          ) : posts.length === 0 ? (
            <Typography>No posts found.</Typography>
          ) : (
            posts.map((post) => (
              <Paper key={post.id} elevation={2} sx={{ p: 3, borderRadius: 3 }}>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  {post.title}
                </Typography>
                <Typography variant="subtitle2" sx={{ color: "#64748b", mb: 2 }}>
                  {post.description}
                </Typography>
                <Button
                  variant="contained"
                  onClick={() => handleViewPost(post.id, post.views + 1, post.likes, post.dislikes)}
                  sx={{
                    background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                    color: "white",
                    borderRadius: 2,
                    fontWeight: 600,
                    py: 1,
                    "&:hover": {
                      transform: "translateY(-2px)",
                      boxShadow: "0 8px 25px rgba(99,102,241,0.6)",
                    },
                  }}
                >
                  Read More
                </Button>
              </Paper>
            ))
          )}
          {loading && <Typography>Loading more posts...</Typography>}
          {!hasMore && <Typography sx={{ textAlign: "center" }}>No more posts.</Typography>}
        </Stack>
      </Container>
    </ThemeProvider>
  );
}
