import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";

//async action to fetch session cookie with the backend
export const fetchSession = createAsyncThunk(
  "auth/fetchSession",
  async (_, { rejectWithValue }) => {
    try {
      const res = await fetch("http://localhost:8080/public/auth/loginStatus", {
        credentials: "include",
      });
      const data = await res.json();
      if (res.ok && data.logged_in) return data;
      return rejectWithValue(null);
    } catch {
      return rejectWithValue(null);
    }
  }
);

interface AuthState {
  isLoggedIn: boolean;
  user: string;
  userID: number;
  loading: boolean;
}

const initialState: AuthState = {
  isLoggedIn: false,
  user: "",
  userID: -1,
  loading: true,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    // setLogin: (state, action: PayloadAction<{ username: string; userID: number }>) => {
    //   state.isLoggedIn = true;
    //   state.user = action.payload.username;
    //   state.userID = action.payload.userID;
    // },
    setLogout: (state) => {
      state.isLoggedIn = false;
      state.user = "";
      state.userID = -1;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchSession.pending, (state) => {
      state.loading = true;
    });
    builder.addCase(fetchSession.fulfilled, (state, action) => {
      state.isLoggedIn = true;
      state.user = action.payload.username;
      state.userID = action.payload.user_id;
      state.loading = false;
    });
    builder.addCase(fetchSession.rejected, (state) => {
      state.isLoggedIn = false;
      state.user = "";
      state.userID = -1;
      state.loading = false;
    });
  },
});

export const { setLogout } = authSlice.actions;
export default authSlice.reducer;