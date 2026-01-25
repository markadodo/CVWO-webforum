"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

import { useState, useEffect, useCallback, useRef } from "react";
import { useSelector } from "react-redux";
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
  Container,
  Typography,
  Box,
  Button,
  TextField,
  Card,
  CardContent,
  CardActions,
  Divider,
  Stack,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  MenuItem,
  FormControl,
  Select,
} from "@mui/material";
import { useRouter, useParams } from "next/navigation";

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
  created_by: number;
  created_at: string;
  username: string;
};

type Post = {
  id: number;
  title: string;
  description: string;
  likes: number;
  dislikes: number;
  views: number;
  popularity: number;
};

export default function TopicDetailPage() {
  const router = useRouter();
  const isLoggedIn = useSelector((state: any) => state.auth.isLoggedIn);
  const user = useSelector((state: any) => state.auth.userID);

  const { topicID } = useParams();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [posts, setPosts] = useState<Post[]>([]);
  const [loadingPosts, setLoadingPosts] = useState(false);
  const [hasMorePosts, setHasMorePosts] = useState(true);
  const [postsPage, setPostsPage] = useState(1);
  const [openCreatePost, setOpenCreatePost] = useState(false);
  const [newPostTitle, setNewPostTitle] = useState("");
  const [newPostDescription, setNewPostDescription] = useState("");
  const [creatingPost, setCreatingPost] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [query, setQuery] = useState("");
  const [isResetting, setIsResetting] = useState(false);
  const [sortBy, setSortBy] = useState<"popularity" | "recency" | "views">("popularity");

  const containerRef = useRef<HTMLDivElement>(null);

  const fetchTopic = async () => {
    try {
      const res = await fetch(`http://localhost:8080/public/topics/${topicID}`, {
        method: "GET",
        credentials: "include"
      });
      const json = await res.json();
      const userRes = await fetch(
        `http://localhost:8080/public/users/${json.created_by}`,
        { method: "GET", credentials: "include" }
      );
      const userJson = await userRes.json();
      setTopic({
          ...json, 
          username: userJson.username
      });
      setTitle(json.title);
      setDescription(json.description);
    } catch (err) {
      console.error(err);
    }
  };
  // Fetch topic
  useEffect(() => {
    fetchTopic();
  }, [topicID]);

  // Edit topic
  const handleSave = async () => {
    try {
      const res = await fetch(`http://localhost:8080/logged_in/topics/${topicID}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ 
          title: title, 
          description: description 
        }),
        credentials: "include",
      });
      const json = await res.json();
      setTopic(json);
      setIsEditing(false);
      fetchTopic();
    } catch (err) {
      console.error(err);
    }
  };

  // Delete topic
  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this topic?")) return;
    try {
      await fetch(`http://localhost:8080/logged_in/topics/${topicID}`, {
        method: "DELETE",
        credentials: "include",
      });
      router.push("/topics");
    } catch (err) {
      console.error(err);
    }
  };

  // Fetch posts
  const fetchPosts = async () => {
    if (!hasMorePosts || loadingPosts) return;

    setLoadingPosts(true);

    try {
      const path = query ? `/search?page=${postsPage}&q=${query}&sort_by=${sortBy}`: `?page=${postsPage}&sort_by=${sortBy}`;
      const res = await fetch(
        `http://localhost:8080/public/topics/${topicID}/posts${path}`
      );
      const json = await res.json();

      if (!json.posts || json.posts.length === 0) {
        setHasMorePosts(false);
        return;
      }

      setPosts(prev => [...prev, ...json.posts]);

      if (json.posts.length < 10) {
        setHasMorePosts(false);
      }

    } catch (err) {
      console.error(err);
    } finally {
      setLoadingPosts(false);
    }
  };

  const handleScroll = useCallback(() => {
    if (!containerRef.current || loadingPosts || !hasMorePosts) return;

    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    if (scrollHeight - scrollTop <= clientHeight + 100) {
      setPostsPage((prev) => prev + 1);
    }
  }, [loadingPosts, hasMorePosts]);

  //Create post
  const handleCreatePost = async () => {
    if (!newPostTitle.trim()) return;

    setCreatingPost(true);
    try {
      const res = await fetch(`http://localhost:8080/logged_in/topics/${topicID}/posts`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          title: newPostTitle,
          description: newPostDescription,
          created_by: user,
        }),
        credentials: "include",
      });

      const json = await res.json();

      setNewPostTitle("");
      setNewPostDescription("");
      setOpenCreatePost(false);

      setPosts((prev) => [json, ...prev]);
      setPostsPage(1);
      setHasMorePosts(true);
    } catch (err) {
      console.error(err);
    } finally {
      setCreatingPost(false);
    }
  };

  useEffect(() => {
    setIsResetting(true);
    setPosts([]);
    setPostsPage(1);
    setHasMorePosts(true);
  }, [topicID, query, sortBy]);

  useEffect(() => {
    if (isResetting) {
      fetchPosts();
      setIsResetting(false);
    } 
    else if(postsPage !== 1) {
      fetchPosts();
    }
  }, [postsPage, isResetting]);

  const handleViewPost = async (postID: number, views: number, likes: number, dislikes: number) => {
    try {
      await fetch(`http://localhost:8080/public/posts/${postID}`, {
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

  if (!topic) return <Typography>Loading...</Typography>;

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Container
        maxWidth="md"
        sx={{ mt: 4, height: '80vh', overflowY: 'auto' }}
        ref={containerRef}
        onScroll={handleScroll}
      >
        {/* Topic Card */}
        <Card sx={{ borderRadius: 3, boxShadow: 3 }}>
          <CardContent>
            {/* Title */}
            {isEditing ? (
              <TextField
                label="Title"
                fullWidth
                sx={{ mb: 2 }}
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            ) : (
              <Typography variant="h4" sx={{ fontWeight: 700, mb: 1 }}>
                {topic.title}
              </Typography>
            )}

            {/* Meta Info */}
            <Typography variant="caption" sx={{ color: "#94a3b8" }}>
              Created by {topic.username} •{" "}
              {new Date(topic.created_at).toLocaleString()}
            </Typography>

            <Divider sx={{ my: 2 }} />

            {/* Description */}
            {isEditing ? (
              <TextField
                label="Description"
                fullWidth
                multiline
                rows={6}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            ) : (
              <Typography variant="body1" sx={{ whiteSpace: "pre-line" }}>
                {topic.description}
              </Typography>
            )}
          </CardContent>

          {/* Actions */}
          {isLoggedIn && topic.created_by === user &&(
            <CardActions sx={{ justifyContent: "flex-end", p: 2 }}>
              {isEditing ? (
                <Stack direction="row" spacing={2}>
                  <Button
                    variant="contained"
                    onClick={handleSave}
                    sx={{
                      background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                    }}
                  >
                    Save
                  </Button>
                  <Button onClick={() => setIsEditing(false)}>Cancel</Button>
                </Stack>
              ) : (
                <Stack direction="row" spacing={2}>
                  <Button
                    variant="outlined"
                    color="primary"
                    onClick={() => setIsEditing(true)}
                  >
                    Edit
                  </Button>
                  <Button variant="outlined" color="error" onClick={handleDelete}>
                    Delete
                  </Button>
                </Stack>
              )}
            </CardActions>
          )}
        </Card>

        {/* Posts under this topic */}
        <Stack>
          <Box   
            display="flex"
            justifyContent="space-between"
            alignItems="center"
            mb={2}
          >
          <form
            onSubmit={(e) => {
              e.preventDefault();
              setQuery(searchTerm);
            }}
          >
            <TextField
              size="small"
              placeholder={`Search posts...`}
              sx={{mt: 1, ml: 1, background: "white", borderRadius: 1, minWidth: 240}}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
            <button type="submit" hidden />
          </form>
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
            {isLoggedIn && (
              <Button
                variant="contained"
                sx={{
                  background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                  color: "white",
                }}
                onClick={() => setOpenCreatePost(true)}
              >
                Create Post
              </Button>
            )}
          </Box>

          <Stack spacing={2}>
            {posts.map((post) => (
              <Card key={post.id} sx={{ borderRadius: 3, boxShadow: 2 }}>
                <CardContent>
                  <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                    {post.title}
                  </Typography>
                  <Typography variant="body2" sx={{ color: "#64748b", mb: 1 }}>
                    {post.description}
                  </Typography>
                  <Typography variant="caption" sx={{ color: "#94a3b8" }}>
                    Views: {post.views} • Likes: {post.likes} • Dislikes: {post.dislikes}
                  </Typography>
                </CardContent>
                <CardActions sx={{ justifyContent: "flex-end", p: 1 }}>
                  <Button
                    variant="contained"
                    size="small"
                    sx={{
                      background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                      color: "white",
                    }}
                    onClick={() =>
                      handleViewPost(post.id, post.views + 1, post.likes, post.dislikes)
                    }
                  >
                    Read More
                  </Button>
                </CardActions>
              </Card>
            ))}

            {posts.length === 0 && !loadingPosts && (
              <Typography>No posts in this topic yet.</Typography>
            )}

            {loadingPosts && <Typography>Loading posts...</Typography>}

            {!hasMorePosts && posts.length > 0 && (
              <Typography sx={{ textAlign: "center" }}>No more posts.</Typography>
            )}
          </Stack>

        </Stack>
        <Dialog open={openCreatePost} onClose={() => setOpenCreatePost(false)}>
          <DialogTitle>Create a New Post</DialogTitle>

          <DialogContent>
            <TextField
              label="Title"
              fullWidth
              sx={{ my: 1 }}
              value={newPostTitle}
              onChange={(e) => setNewPostTitle(e.target.value)}
            />

            <TextField
              label="Description"
              fullWidth
              multiline
              rows={4}
              sx={{ my: 1 }}
              value={newPostDescription}
              onChange={(e) => setNewPostDescription(e.target.value)}
            />
          </DialogContent>

          <DialogActions>
            <Button onClick={() => setOpenCreatePost(false)}>
              Cancel
            </Button>

            <Button
              variant="contained"
              disabled={creatingPost || !newPostTitle.trim()}
              onClick={handleCreatePost}
              sx={{
                background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                color: "white",
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