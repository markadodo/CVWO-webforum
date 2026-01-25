"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

import { useState, useEffect, useRef, useCallback } from "react";
import { useSelector } from "react-redux";
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
  Container,
  Typography,
  Stack,
  Paper,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
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

type Topic = {
  id: number;
  title: string;
  description: string;
  created_by: string;
  created_at: string;
};

export default function TopicsPage() {
  const isLoggedIn = useSelector((state: any) => state.auth.isLoggedIn);
  const router = useRouter();
  const [query, setQuery] = useState<string | null>(null);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    setQuery(params.get("q"));
  }, []);

  const [topics, setTopics] = useState<Topic[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const [openCreate, setOpenCreate] = useState(false);
  const [newTitle, setNewTitle] = useState("");
  const [newDescription, setNewDescription] = useState("");

  const containerRef = useRef<HTMLDivElement>(null);
  const [isResetting, setIsResetting] = useState(false);

  // Fetch topics
  const fetchTopics = async () => {
    if (!hasMore || loading) return;

    setLoading(true);
    try {
      const path = query ? `/search?page=${page}&q=${query}` : `?page=${page}`;

      const res = await fetch(`${apiUrl}/public/topics${path}`);
      const json = await res.json();

      if (json.topics.length === 0) {
        setHasMore(false);
        return;
      }

      setTopics((prev) => [...prev, ...json.topics]);
      if (json.topics.length < 10) setHasMore(false);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    setIsResetting(true);
    setTopics([]);
    setPage(1);
    setHasMore(true);
  }, [query]);

  useEffect(() => {
    if (isResetting) {
      fetchTopics();
      setIsResetting(false);
    } else if(page !== 1) {
      fetchTopics();
    }
  }, [page, isResetting]);

  // Infinite scroll
  const handleScroll = useCallback(() => {
    if (!containerRef.current || loading || !hasMore) return;

    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    if (scrollHeight - scrollTop <= clientHeight + 100) {
      setPage((prev) => prev + 1);
    }
  }, [loading, hasMore]);

  // Create topic
  const handleCreateTopic = async () => {
    if (!newTitle || !newDescription) return;

    try {
      const res = await fetch(`${apiUrl}/logged_in/topics`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          title: newTitle,
          description: newDescription,
        }),
        credentials: "include",
      });

      const json = await res.json();
      setTopics((prev) => [json, ...prev]);
      setOpenCreate(false);
      setNewTitle("");
      setNewDescription("");
    } catch (err) {
      console.error("Error creating topic:", err);
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
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h4" sx={{ fontWeight: 700 }}>
            All Topics
          </Typography>
          <Button
            variant="contained"
            onClick={() => setOpenCreate(true)}
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
            Create Topic
          </Button>
        </Box>

        <Stack spacing={3}>
          {topics.length === 0 && loading ? (
            <Typography>Loading topics...</Typography>
          ) : topics.length === 0 ? (
            <Typography>No topics found.</Typography>
          ) : (
            topics.map((topic) => (
              <Paper key={topic.id} elevation={2} sx={{ p: 3, borderRadius: 3 }}>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  {topic.title}
                </Typography>
                <Typography variant="subtitle2" sx={{ color: "#64748b", mb: 1 }}>
                  {topic.description}
                </Typography>
                <Button
                  variant="contained"
                  onClick={() => router.push(`/topics/${topic.id}`)}
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
          {loading && <Typography>Loading more topics...</Typography>}
          {!hasMore && <Typography sx={{ textAlign: "center" }}>No more topics.</Typography>}
        </Stack>

        {/* Create Topic Dialog */}
        <Dialog open={openCreate} onClose={() => setOpenCreate(false)}>
          <DialogTitle>Create a New Topic</DialogTitle>
          <DialogContent>
            {!isLoggedIn && (
              <Typography color="error" sx={{ mb: 2 }}>
                You must be logged in to create a topic.
              </Typography>
            )}
            <TextField
              label="Title"
              fullWidth
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
              sx={{ my: 1 }}
              disabled={!isLoggedIn}
            />
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={4}
              value={newDescription}
              onChange={(e) => setNewDescription(e.target.value)}
              sx={{ my: 1 }}
              disabled={!isLoggedIn}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpenCreate(false)}>Cancel</Button>
            <Button
              onClick={handleCreateTopic}
              variant="contained"
              disabled={!isLoggedIn || !newTitle || !newDescription}
              sx={{
                background: !isLoggedIn
                  ? "#ccc"
                  : "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                color: !isLoggedIn ? "#666" : "white",
                "&:hover": !isLoggedIn
                  ? {}
                  : {
                      transform: "translateY(-2px)",
                      boxShadow: "0 8px 25px rgba(99,102,241,0.6)",
                    },
              }}
            >
              Create
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </ThemeProvider>
  );
}
