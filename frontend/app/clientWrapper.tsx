"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

import React, { ReactNode, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import type { RootState, AppDispatch } from "@/lib/store";
import { setLogout, fetchSession } from "@/lib/features/authSlice";
import { useRouter } from "next/navigation";

import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  CssBaseline,
  Container,
  TextField,
  Button,
  Drawer,
  List,
  ListItemText,
  Divider,
  IconButton,
  ListItemButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Stack,
  Alert,
  Menu,
  MenuItem,
  Select,
  SelectChangeEvent,
  FormControl,
  ListItemIcon,
  createTheme,
  ThemeProvider
} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import WhatshotRoundedIcon from "@mui/icons-material/WhatshotRounded";
import ScheduleRoundedIcon from "@mui/icons-material/ScheduleRounded";
import VisibilityRoundedIcon from "@mui/icons-material/VisibilityRounded";

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

const drawerWidth = 240;
const collapsedWidth = 64;

type InputFields = "username" | "password";
type Errors = {
  username?: string;
  password?: string;
  login?: string;
  register?: string;
};

const sortOptions = [
  {
    label: "Popular",
    value: "popularity",
    color: "#F57C00",
    icon: <WhatshotRoundedIcon />,
  },
  {
    label: "Recent",
    value: "recency",
    color: "#1976D2",
    icon: <ScheduleRoundedIcon />,
  },
  {
    label: "Views",
    value: "views",
    color: "#7B1FA2",
    icon: <VisibilityRoundedIcon />,
  },
];

export default function ClientLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const dispatch = useDispatch<AppDispatch>();
  const { isLoggedIn, user, userID } = useSelector((state: RootState) => state.auth);

  const [open, setOpen] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchType, setSearchType] = useState<"posts" | "topics">("posts");
  const [openAuth, setOpenAuth] = useState(false);
  const [inputInfo, setInputInfo] = useState({ username: "", password: "" });
  const [error, setError] = useState<Errors>({});
  const [successMsg, setSuccessMsg] = useState<string>("");
  const [authMode, setAuthMode] = useState<"login" | "register">("login");

  const [anchorElUser, setAnchorElUser] = useState<null | HTMLElement>(null);


  useEffect(() => {
    dispatch(fetchSession());
  }, [dispatch]);


  //Close auth dialog automatically when logged in
  useEffect(() => {
    if (isLoggedIn) setOpenAuth(false);
  }, [isLoggedIn]);

  
  async function loginHandler(username: string, password: string) {
    const res = await fetch("http://localhost:8080/public/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
      credentials: "include"
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Failed to login");
    return data;
  }

  async function registerHandler(username: string, password: string) {
    const res = await fetch("http://localhost:8080/public/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Failed to register");
    return data;
  }

  // Login function 
  async function login(username: string, password: string) {
    try {
      const data = await loginHandler(username, password);

      // await new Promise(r => setTimeout(r, 5000));
      await dispatch(fetchSession()).unwrap();

      setInputInfo({ username: "", password: "" });
      setError({});
      setSuccessMsg("Logged in successfully!");
    } catch (err: any) {
      setError({ login: err.message });
      setSuccessMsg("");
    }
  }

  // Register function 
  async function register(username: string, password: string) {
    try {
      const data = await registerHandler(username, password);

      setInputInfo({ username: "", password: "" });
      setError({});
      setSuccessMsg("Registration successful! You can now log in.");
      setAuthMode("login");
    } catch (err: any) {
      setError({ register: err.message });
      setSuccessMsg("");
    }
  }

  // Input change handler 
  const handleInputInfoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const key = e.target.name as InputFields;
    setInputInfo({ ...inputInfo, [key]: e.target.value });
    setError({});
    setSuccessMsg("");
  };

  // Form submission 
  const handleSubmit = async () => {
    const newError: Errors = {};
    if (!inputInfo.username) newError.username = "Username cannot be empty";
    if (!inputInfo.password) newError.password = "Password cannot be empty";
    if (Object.keys(newError).length > 0) {
      setError(newError);
      return;
    }

    if (authMode === "login") await login(inputInfo.username, inputInfo.password);
    else await register(inputInfo.username, inputInfo.password);
  };

  const handleLogout = async () => {
    await fetch("http://localhost:8080/public/auth/logout", {
      method: "POST",
      credentials: "include",
    });

    dispatch(setLogout());
    setSuccessMsg("");
  };

  const handleSearchTypeChange = (event: SelectChangeEvent<"posts" | "topics">) => {
    setSearchType(event.target.value as "posts" | "topics");
  };

  const handleUserMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleUserMenuClose = (route?: string) => {
    setAnchorElUser(null);
    if (route) router.push(route);
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />

      {/* AUTH DIALOG */}
      <Dialog open={openAuth} onClose={() => setOpenAuth(false)}>
        <DialogTitle>{authMode === "login" ? "Login to ChatIt" : "Register for ChatIt"}</DialogTitle>
        <DialogContent>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handleSubmit();
            }}
          >
            <Stack spacing={2} mt={2}>
              <TextField
                size="small"
                name="username"
                placeholder="username"
                value={inputInfo.username}
                onChange={(e) => setInputInfo({ ...inputInfo, username: e.target.value })}
                error={!!error.username}
                helperText={error.username}
              />
              <TextField
                size="small"
                name="password"
                type="password"
                placeholder="password"
                value={inputInfo.password}
                onChange={(e) => setInputInfo({ ...inputInfo, password: e.target.value })}
                error={!!error.password}
                helperText={error.password}
              />

              <Button type="submit" variant="contained" sx={{
                background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                color: "white", borderRadius: 2, fontWeight: 600, py: 1,
                "&:hover": { transform: "translateY(-2px)", boxShadow: "0 8px 25px rgba(99,102,241,0.6)" },
              }}>
                {authMode === "login" ? "Login" : "Register"}
              </Button>

              {!!error.login && <Alert severity="error">{error.login}</Alert>}
              {!!error.register && <Alert severity="error">{error.register}</Alert>}
              {!!successMsg && <Alert severity="success">{successMsg}</Alert>}

              <Button onClick={() => setAuthMode(authMode === "login" ? "register" : "login")} sx={{ textTransform: "none", mt: 1 }}>
                {authMode === "login" ? "Don't have an account? Register" : "Already have an account? Login"}
              </Button>
            </Stack>
          </form>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAuth(false)}>Cancel</Button>
        </DialogActions>
      </Dialog>

      {/* APP BAR */}
      <AppBar position="fixed" sx={{ boxShadow: 0 }}>
        <Toolbar>
          <IconButton color="inherit" edge="start" onClick={() => setOpen(!open)}>
            <MenuIcon />
          </IconButton>

          <Typography variant="h6" sx={{ ml: 2 }}>ChatIt</Typography>

          {/* SEARCH WITH DROPDOWN */}
          <form
            onSubmit={(e) => {
              e.preventDefault();
              router.push(`/${searchType}?q=${encodeURIComponent(searchTerm)}`);
            }}
          >
            <FormControl size="small" sx={{ ml: 3, background: "white", borderRadius: 1, minWidth: 60 }}>
              <Select value={searchType} onChange={handleSearchTypeChange}>
                <MenuItem value="posts">Posts</MenuItem>
                <MenuItem value="topics">Topics</MenuItem>
              </Select>
            </FormControl>

            <TextField
              size="small"
              placeholder={`Search ${searchType}...`}
              sx={{ ml: 1, background: "white", borderRadius: 1, minWidth: 240}}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />

            <button type="submit" hidden />
          </form>

          <Box sx={{ ml: "auto" }}>
            {isLoggedIn ? (
              <Stack direction="row" spacing={2} alignItems="center">
                <Button color="inherit" onClick={handleUserMenuOpen} endIcon={<ArrowDropDownIcon />}>
                  {user}
                </Button>
                <Menu
                  anchorEl={anchorElUser}
                  open={Boolean(anchorElUser)}
                  onClose={() => handleUserMenuClose()}
                >
                  <MenuItem onClick={() => handleUserMenuClose("/profile")}>Profile</MenuItem>
                  <MenuItem onClick={() => { handleLogout(); handleUserMenuClose(); }}>Logout</MenuItem>
                </Menu>
              </Stack>
            ) : (
              <Button color="inherit" onClick={() => setOpenAuth(true)}>Login / Register</Button>
            )}
          </Box>
        </Toolbar>
      </AppBar>

      {/* SIDEBAR */}
      <Drawer
        variant="permanent"
        sx={{
          width: open ? drawerWidth : collapsedWidth,
          [`& .MuiDrawer-paper`]: { width: open ? drawerWidth : collapsedWidth, mt: 8, transition: "width 0.2s" },
        }}
      >
        <List sx={{ p: 2, opacity: open ? 1 : 0 }}>
          <Typography variant="caption" sx={{ mb: 1 }}>Posts</Typography>
          {sortOptions.map(({ label, value, color, icon }) => (
            <ListItemButton
              key={value}
              onClick={() => router.push(`/posts?sort_by=${value}`)}
              sx={{
                borderRadius: 1,
                "& .MuiListItemIcon-root": {
                  minWidth: 36,
                  color: color,
                },
                "&:hover": {
                  bgcolor: `${color}14`,
                },
              }}
            >
              <ListItemIcon>
                {icon}
              </ListItemIcon>

              <ListItemText primary={label} />
            </ListItemButton>
          ))}

            {/* Divider */}
          <Divider sx={{ my: 2 }} />

          {/* Topics section */}
          <Typography variant="caption" sx={{ mb: 1 }}>
            Topics
          </Typography>

          <ListItemButton
            key={"topic"}
            onClick={() => router.push(`/topics`)}
            sx={{
              borderRadius: 1,
              "& .MuiListItemIcon-root": {
                minWidth: 36,
                color: "black",
              },
              "&:hover": {
                bgcolor: `${"grey"}14`,
              },
            }}
          >
            <ListItemText primary="Main Page" />
          </ListItemButton>
        </List>
      </Drawer>

      {/* MAIN CONTENT */}
      <Box component="main" sx={{ ml: open ? `${drawerWidth}px` : `${collapsedWidth}px`, mt: 10, transition: "margin-left 0.3s" }}>
        <Container maxWidth="lg">{children}</Container>
      </Box>
    </ThemeProvider>
  );
}
