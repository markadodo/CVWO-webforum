"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

import { useSelector } from "react-redux";
import React, { useEffect, useState, useRef, useCallback } from "react";
import {
  Box,
  Container,
  Typography,
  Paper,
  Divider,
  Stack,
  Button,
  IconButton,
  TextField,
  ThemeProvider,
  createTheme,
  CssBaseline,
  CardActions,
  Card,
} from "@mui/material";
import { ThumbUp, ThumbDown } from "@mui/icons-material";
import WhatshotRoundedIcon from "@mui/icons-material/WhatshotRounded";
import AddIcon from "@mui/icons-material/Add";
import { useParams, useRouter } from "next/navigation";

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
  topic_id: number;
  likes: number;
  dislikes: number;
  is_edited: boolean;
  views: number;
  popularity: number;
  created_by: number;
  created_at: string;
  username: string;
};

type Comment = {
  id: number;
  description: string;
  likes: number;
  dislikes: number;
  is_edited: boolean;
  created_by: string;
  created_at: string;
  username: string;
  replies?: Comment[];
  replyPage?: number;
  hasMoreReplies?: boolean;
  loadingReplies?: boolean;
};

export default function PostPage() {
  const isLoggedIn = useSelector((state: any) => state.auth.isLoggedIn);
  const user = useSelector((state: any) => state.auth.userID);
  const usern = useSelector((state: any) => state.auth.user);
  const params = useParams();
  const postID: number = params?.postID
      ? Array.isArray(params.postID)
      ? parseInt(params.postID[0]) 
      : parseInt(params.postID)   
      : 0;                           

  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [userReaction, setUserReaction] = useState<"like" | "dislike" | null>(null);

  const [comments, setComments] = useState<Comment[]>([]);
  const [commentsPage, setCommentsPage] = useState(1);
  const [hasMoreComments, setHasMoreComments] = useState(true);
  const [loadingComments, setLoadingComments] = useState(false);
  const [isResettingComments, setIsResettingComments] = useState(false);
  
  const [replyingToId, setReplyingToId] = useState<number | null>(null);
  const [replyText, setReplyText] = useState("");

  const [isEditing, setIsEditing] = useState(false);
  const [editTitle, setEditTitle] = useState(post?.title || "");
  const [editDescription, setEditDescription] = useState(post?.description || "");


  const router = useRouter();
  const containerRef = useRef<HTMLDivElement>(null);

  const insertReply = (
    comments: Comment[],
    parentId: number,
    newReply: Comment
  ): Comment[] => {
    return comments.map((c) => {
      if (c.id === parentId) {
        return {
          ...c,
          replies: [...(c.replies || []), newReply],
        };
      }

      if (c.replies?.length) {
        return {
          ...c,
          replies: insertReply(c.replies, parentId, newReply),
        };
      }

      return c;
    });
  };

  const updateCommentInTree = (
    comments: Comment[],
    targetId: number,
    updater: (c: Comment) => Comment
  ): Comment[] => {
    return comments.map((c) => {
      if (c.id === targetId) return updater(c);
      if (c.replies?.length) {
        return { ...c, replies: updateCommentInTree(c.replies, targetId, updater) };
      }
      return c;
    });
  };

  const submitComment = async () => {
    if (!replyText.trim()) return;

    try {
      const res = await fetch(`${apiUrl}/logged_in/comments`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          post_id: postID,
          description: replyText,
        }),
        credentials: "include",
      });

      const newComment = await res.json();

      const newCommentWithFields = {
        ...newComment,
        replies: [],
        replyPage: 1,
        hasMoreReplies: true,
        loadingReplies: false,
      };

      setComments((prev) => [newCommentWithFields, ...prev]);
      setReplyText("");
    } catch (err) {
      console.error(err);
    }
  };


  const submitReply = async (parentID: number) => {
    if (!replyText.trim()) return;

    const res = await fetch(`${apiUrl}/logged_in/comments`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        post_id: postID,
        parent_comment_id: parentID,
        description: replyText,
      }),
      credentials: "include",
    });

     const newReply = await res.json();

    const newReplyWithFields = {
      ...newReply,
      username: usern,
      replies: [],
      replyPage: 1,
      hasMoreReplies: true,
      loadingReplies: false,
    };

    setComments((prev) => insertReply(prev, parentID, newReplyWithFields));
    setReplyText("");
    setReplyingToId(null);
  };


  function CommentItem({
    comment,
    depth = 0,
    onViewReplies,
  }: {
    comment: Comment;
    depth?: number;
    onViewReplies: (c: Comment) => void;
  }) {
    return (
      <Box
        sx={{
          pl: depth ? 2 : 0,
          ml: depth ? 1 : 0,
          borderLeft: depth ? "2px solid #e5e7eb" : "none",
          mt: 1.5,
        }}
      >
        <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 0.5 }}>
          <Typography sx={{ fontWeight: 600, fontSize: 14, zIndex: 1 }}>
            {comment.username}
          </Typography>

          <Divider flexItem sx={{ borderBottomWidth: 2, backgroundColor: "#e5e7eb", mx: 1 }} />

          {isLoggedIn && comment.created_by === user && comment.description !== "Comment deleted" && (
            <Button
              size="small"
              color="error"
              sx={{ textTransform: "none", minWidth: 50, p: "2px 6px", zIndex: 1 }}
              onClick={() => handleDeleteComment(comment.id)}
            >
              Delete
            </Button>
          )}
        </Stack>

        <Typography
          sx={{
            fontSize: 14,
            color: comment.description === "Comment deleted" ? "error.main" : "text.primary",
          }}
        >
          {comment.description}
        </Typography>

        <Stack direction="row" spacing={1} sx={{ mt: 0.5 }}>
          {isLoggedIn && (<Button
            size="small"
            sx={{ textTransform: "none", px: 0 }}
            onClick={() => setReplyingToId(comment.id)}
          >
            Reply
          </Button>)}

          {comment.hasMoreReplies && !comment.loadingReplies &&(
            <Button
              size="small"
              sx={{ textTransform: "none", px: 0 }}
              onClick={() => onViewReplies(comment)}
            >
              View replies
            </Button>
          )}
        </Stack>

        {replyingToId === comment.id && (
          <Box sx={{ mt: 1, opacity: isLoggedIn ? 1 : 0.5, pointerEvents: isLoggedIn ? "auto" : "none" }}>
            <TextField
              fullWidth
              size="small"
              autoFocus
              placeholder="Write a reply..."
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
            />
            {isLoggedIn && (<Stack direction="row" spacing={1} sx={{ mt: 0.5 }}>
              <Button size="small" onClick={() => submitReply(comment.id)}>
                Post
              </Button>
              <Button
                size="small"
                color="inherit"
                onClick={() => setReplyingToId(null)}
              >
                Cancel
              </Button>
            </Stack>)}
          </Box>
        )}

        {comment.replies?.map((reply) => (
          <CommentItem
            key={reply.id}
            comment={reply}
            depth={depth + 1}
            onViewReplies={onViewReplies}
          />
        ))}
      </Box>
    );
  }
  
  // Fetch post
  useEffect(() => {
    if (!postID) return;
    const fetchPost = async () => {
      setLoading(true);
      try {
        const res = await fetch(`${apiUrl}/public/posts/${postID}`);
        const postJson = await res.json();

        const userRes = await fetch(
          `${apiUrl}/public/users/${postJson.created_by}`,
          { method: "GET", credentials: "include" }
        );
        const userJson = await userRes.json();

        setPost({
          ...postJson,
          username: userJson.username,
        });
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchPost();
  }, [postID]);

  // Fetch comments
  const fetchComments = async () => {
    if (!hasMoreComments || loadingComments) return;
    setLoadingComments(true);
    try {
      const res = await fetch(
        `${apiUrl}/public/posts/${postID}/comments?page=${commentsPage}&order=DESC`
      );
      const json = await res.json();
      if (json.comments.length === 0) {
        setHasMoreComments(false);
        return;
      }

      const newComments = await Promise.all(
        (json.comments || []).map(async (c: Comment) => {
          const userRes = await fetch(`${apiUrl}/public/users/${c.created_by}`, {
            credentials: "include",
          });
          const userJson = await userRes.json();
          return {
            ...c,
            username: userJson.username,
            description: c.description === "" ? "Comment deleted" : c.description,
            replies: [],
            replyPage: 1,
            hasMoreReplies: true,
            loadingReplies: false,
          };
        })
      );
      setComments((prev) => [...prev, ...newComments]);
      if (json.comments.length < 10) setHasMoreComments(false);
    } catch (err) {
      console.error(err);
    } finally {
      setLoadingComments(false);
    }
  };

  useEffect(() => {
    setIsResettingComments(true);
    setComments([]);
    setCommentsPage(1);
    setHasMoreComments(true);
  }, [postID]);

  useEffect(() => {
    if (isResettingComments) {
      fetchComments();
      setIsResettingComments(false);
    } else if (commentsPage !== 1) {
      fetchComments();
    }
  }, [commentsPage, isResettingComments]);

  const handleCommentsScroll = useCallback(() => {
    if (!containerRef.current || loadingComments || !hasMoreComments) return;
    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    if (scrollHeight - scrollTop <= clientHeight + 100) {
      setCommentsPage((prev) => prev + 1);
    }
  }, [loadingComments, hasMoreComments]);

  // REPLACE handleViewReplies WITH THIS
  const handleViewReplies = async (parent: Comment) => {
    if (!parent.hasMoreReplies || parent.loadingReplies) return;

    setComments((prev) =>
      updateCommentInTree(prev, parent.id, (c) => ({
        ...c,
        loadingReplies: true,
      }))
    );

    try {
      const res = await fetch(
        `${apiUrl}/public/comments/${parent.id}?page=${parent.replyPage}`
      );
      const json = await res.json();
      const newReplies = await Promise.all(
        (json.comments || []).map(async (c: Comment) => {
          const userRes = await fetch(`${apiUrl}/public/users/${c.created_by}`, {
            credentials: "include",
          });
          const userJson = await userRes.json();
          return {
            ...c,
            username: userJson.username,
            description: c.description === "" ? "Comment deleted" : c.description,
            replies: [],
            replyPage: 1,
            hasMoreReplies: true,
            loadingReplies: false,
          };
        })
      );

      setComments(prev =>
        updateCommentInTree(prev, parent.id, c => {
          const existingIds = new Set(c.replies?.map(r => r.id));
          const filteredNewReplies = newReplies.filter(r => !existingIds.has(r.id));
          return {
            ...c,
            replies: [...(c.replies || []), ...filteredNewReplies],
            replyPage: (c.replyPage ?? 1) + 1,
            hasMoreReplies: filteredNewReplies.length === 10,
            loadingReplies: false,
          };
        })
      );
    } catch (err) {
      console.error(err);
      setComments(prev =>
        updateCommentInTree(prev, parent.id, c => ({ ...c, loadingReplies: false }))
      );
    }
  };

    // Fetch current user reaction
  useEffect(() => {
    const fetchReaction = async () => {
      if (!post) return;
      try {
        const res = await fetch(`${apiUrl}/logged_in/posts/${post.id}/reactions`, {
          method: "GET",
          credentials: "include",
        });
        const json = await res.json();

        if (json.reaction === true) {
          setUserReaction("like");
        } else if (json.reaction === false) {
          setUserReaction("dislike");
        } else if (json.reaction === null) {
          setUserReaction(null);
        }
      } catch (err) {
        console.error(err);
      }
    };

    fetchReaction();
  }, [post]);

  // Handle Like
  const handleLike = async () => {
    if (!post) return;

    if (userReaction === "dislike") {
      alert("You already disliked this post. Remove that first.");
      return;
    }

    try {
      if (userReaction === null) {
        await fetch(`${apiUrl}/logged_in/posts/${post.id}/reactions`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ reaction: true }),
          credentials: "include",
        });
        post.likes++;
        setUserReaction("like");
      } else {
        await fetch(`${apiUrl}/logged_in/posts/${post.id}/reactions`, {
          method: "DELETE",
          credentials: "include",
        });
        post.likes--;
        setUserReaction(null);
      }
    } catch (err) {
      console.error(err);
    }
  };

  // Handle Dislike
  const handleDislike = async () => {
    if (!post) return;

    if (userReaction === "like") {
      alert("You already liked this post. Remove that first.");
      return;
    }

    try {
      if (userReaction === null) {
        await fetch(`${apiUrl}/logged_in/posts/${post.id}/reactions`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ reaction: false }),
          credentials: "include",
        });
        post.dislikes++;
        setUserReaction("dislike");
      } else {
        await fetch(`${apiUrl}/logged_in/posts/${post.id}/reactions`, {
          method: "DELETE",
          credentials: "include",
        });
        post.dislikes--;
        setUserReaction(null);
      }
    } catch (err) {
      console.error(err);
    }
  };

  const handleSavePost = async () => {
    if (!post) return;
    try {
      const res = await fetch(`${apiUrl}/logged_in/posts/${postID}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ 
          title: editTitle, 
          description: editDescription 
        }),
        credentials: "include",
      });
      const json = await res.json();
      setPost(json);
      setIsEditing(false);
    } catch (err) {
      console.error(err);
    }
  };

  const handleDeletePost = async () => {
    if (!confirm("Are you sure you want to delete this post?")) return;
    try {
      await fetch(`${apiUrl}/logged_in/posts/${postID}`, {
        method: "DELETE",
        credentials: "include",
      });
      router.push("/topics");
    } catch (err) {
      console.error(err);
    }
  };

  const handleDeleteComment = async (commentID: number) => {
    if (!confirm("Are you sure you want to delete this post?")) return;
    try {
      await fetch(`${apiUrl}/logged_in/comments/${commentID}`, {
        method: "DELETE",
        credentials: "include",
      });
      setComments((prev) =>
        updateCommentInTree(prev, commentID, (c) => ({
          ...c,
          description: "Comment deleted",
        }))
      );
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) {
    return (
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Container maxWidth="md" sx={{ mt: 4 }}>
          <Typography>Loading post...</Typography>
        </Container>
      </ThemeProvider>
    );
  }

  if (!post) {
    return (
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Container maxWidth="md" sx={{ mt: 4 }}>
          <Typography>Post not found.</Typography>
        </Container>
      </ThemeProvider>
    );
  }

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Container
        maxWidth="md"
        sx={{ mt: 4, height: "80vh", overflowY: "auto" }}
        ref={containerRef}
        onScroll={handleCommentsScroll}
      >

        {/* Post Content */}
        <Card elevation={2} sx={{ p: 3, borderRadius: 3, backgroundColor: "#ffffff", mb: 4 }}>
          {isEditing ? (
            <TextField
              label="Title"
              fullWidth
              sx={{ mb: 2 }}
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
            />
          ) : (
            <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>
              {post.title}
            </Typography>
          )}

          <Typography variant="subtitle2" sx={{ color: "#64748b", mt: 1 }}>
            by {post.username} • {new Date(post.created_at).toLocaleDateString()} • {post.views} views
            {post.is_edited ? " • Edited" : ""}
          </Typography>

          {isEditing ? (
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={4}
              sx={{ mb: 2 }}
              value={editDescription}
              onChange={(e) => setEditDescription(e.target.value)}
            />
          ) : (
            <Typography sx={{ mb: 2 }}>{post.description}</Typography>
          )}
          <Divider sx={{ my: 2 }} />

          {/* Likes / Dislikes / Popularity */}
          <Stack direction="row" spacing={2} alignItems="center" justifyContent="space-between">
            <Stack direction="row" spacing={1} alignItems="center">
              <IconButton 
                color={userReaction === "like" && isLoggedIn ? "primary" : "default"}
                onClick={() => {
                  if (isLoggedIn) {
                    handleLike();
                  }
                  return
                }}
              >
                <ThumbUp />
              </IconButton>
              <Typography>{post.likes}</Typography>

              <IconButton 
                color={userReaction === "dislike" && isLoggedIn ? "secondary" : "default"}
                onClick={() => {
                  if (isLoggedIn) {
                    handleDislike();
                  }
                  return
                }}
              >
                <ThumbDown />
              </IconButton>
              <Typography>{post.dislikes}</Typography>
            </Stack>

            <Stack direction="row" spacing={1} alignItems="center">
              <WhatshotRoundedIcon sx={{ color: "#F57C00" }} />
              <Typography sx={{ ml: 1 }}>{post.popularity}</Typography>
            </Stack>

            {/* Post Actions: Edit / Delete */}
            {isLoggedIn && post.created_by === user && (
              <CardActions sx={{ justifyContent: "flex-end", p: 2 }}>
                {isEditing ? (
                  <Stack direction="row" spacing={2}>
                    <Button
                      variant="contained"
                      onClick={handleSavePost}
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
                    <Button
                      variant="outlined"
                      color="error"
                      onClick={handleDeletePost}
                    >
                      Delete
                    </Button>
                  </Stack>
                )}
              </CardActions>
            )}
          </Stack>
        </Card>

        {/* Comments Section */}

        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            mt: 2,
            mb: 1,
            cursor: isLoggedIn ? "pointer" : "not-allowed", // visual cue
            opacity: isLoggedIn ? 1 : 0.5, // faded if not logged in
          }}
          onClick={() => isLoggedIn && setReplyingToId(null)}
        >
          <AddIcon fontSize="small" sx={{ mr: 1, color: "#6366f1" }} />
          <Typography sx={{ fontSize: 14, fontWeight: 500, color: "#6366f1" }}>
            Add a comment
          </Typography>
        </Box>

        {/* Top-level comment input */}
        {replyingToId === null && (
          <Box
            sx={{
              mb: 2,
              opacity: isLoggedIn ? 1 : 0.5,
              pointerEvents: isLoggedIn ? "auto" : "none", // prevents typing if not logged in
            }}
          >
            <TextField
              fullWidth
              size="small"
              placeholder={isLoggedIn ? "Write a comment..." : "Login to comment"}
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
            />
            {isLoggedIn && (
              <Stack direction="row" spacing={1} sx={{ mt: 0.5 }}>
                <Button size="small" onClick={submitComment}>
                  Post
                </Button>
              </Stack>
            )}
          </Box>
        )}
        {comments.map((comment) => (
          <CommentItem
            key={comment.id}
            comment={comment}
            onViewReplies={handleViewReplies}
          />
        ))}
      </Container>
    </ThemeProvider>
  );
}
