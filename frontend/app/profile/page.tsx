"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

import { useSelector } from "react-redux";
import React, { useEffect, useState } from "react";
import {
  Box,
  Container,
  Typography,
  Avatar,
  Paper,
  Stack,
  Button,
  TextField,
  Alert,
} from "@mui/material";
import { useRouter } from "next/navigation";

export default function ProfilePage() {
  const router = useRouter();
  const user = useSelector((state: any) => state.auth.userID);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const passwordsMatch =
    password === confirmPassword || confirmPassword === "";

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await fetch(
          `http://localhost:8080/logged_in/users/${user}`,
          { credentials: "include" }
        );

        if (!res.ok) throw new Error("Failed to load profile");

        const json = await res.json();
        setUsername(json.username);
      } catch (err) {
        setError("Unable to load profile");
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  const handleSave = async () => {
    if (!passwordsMatch) return;

    setSaving(true);
    setError(null);
    setSuccess(false);

    try {
      const body: any = { username };

      if (password.trim() !== "") {
        body.password = password;
      }

      const res = await fetch(
        `http://localhost:8080/logged_in/users/${user}`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(body),
        }
      );

      if (!res.ok) throw new Error("Update failed");

      setPassword("");
      setConfirmPassword("");
      setSuccess(true);
      window.location.href = "/";
    } catch (err) {
      setError("Failed to save changes");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <Container maxWidth="sm">
        <Typography>Loading profile…</Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm">
      <Paper
        elevation={0}
        sx={{
          p: 4,
          borderRadius: 4,
          border: "1px solid #e2e8f0",
          background: "white",
        }}
      >
        <Stack spacing={3}>
          {/* Header */}
          <Stack direction="row" spacing={2} alignItems="center">
            <Avatar
              sx={{
                width: 72,
                height: 72,
                background: "linear-gradient(135deg, #6366f1, #8b5cf6)",
                fontSize: "2rem",
                fontWeight: 700,
              }}
            >
              {username?.[0]?.toUpperCase()}
            </Avatar>

            <Typography variant="h5" fontWeight={700}>
              Edit Profile
            </Typography>
          </Stack>

          {error && <Alert severity="error">{error}</Alert>}
          {success && <Alert severity="success">Profile updated</Alert>}

          {/* Username */}
          <Box>
            <Typography variant="caption" fontWeight={600}>
              Username
            </Typography>
            <TextField
              fullWidth
              size="small"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </Box>

          {/* Password */}
          <Box>
            <Typography variant="caption" fontWeight={600}>
              New Password
            </Typography>
            <TextField
              fullWidth
              size="small"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </Box>

          {/* Confirm Password */}
          <Box>
            <Typography variant="caption" fontWeight={600}>
              Confirm Password
            </Typography>
            <TextField
              fullWidth
              size="small"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              error={!passwordsMatch}
              helperText={!passwordsMatch ? "Passwords do not match" : ""}
            />
          </Box>

          {/* Save */}
          <Button
            variant="contained"
            disabled={!passwordsMatch || saving}
            onClick={handleSave}
            sx={{
              mt: 1,
              textTransform: "none",
              fontWeight: 600,
              background:
                "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
            }}
          >
            {saving ? "Saving..." : "Save Changes"}
          </Button>
        </Stack>
      </Paper>
    </Container>
  );
}
