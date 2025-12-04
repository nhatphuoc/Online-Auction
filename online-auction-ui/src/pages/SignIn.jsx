import { useState } from "react";
import { signIn } from "../services/authService";

export default function SignIn() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();

    const res = await signIn({ email, password });

    if (!res.success) return alert("Invalid credentials");

    localStorage.setItem("accessToken", res.accessToken);
    localStorage.setItem("refreshToken", res.refreshToken);

    alert("Login success!");
    // chuyển về homepage
  }

  return (
    <div className="container mt-5" style={{ maxWidth: 450 }}>
      <h3 className="text-center mb-4">Sign-In</h3>

      <form onSubmit={handleSubmit}>
        <div className="mb-3">
          <label>Email</label>
          <input className="form-control" value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>

        <div className="mb-3">
          <label>Password</label>
          <input
            type="password"
            className="form-control"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>

        <button className="btn btn-primary w-100">Sign-In</button>
      </form>

      <div className="text-center mt-3">Or</div>

      <button
        className="btn btn-danger w-100 mt-3"
        onClick={() => {
          const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID;

          if (!GOOGLE_CLIENT_ID) {
            console.error("GOOGLE CLIENT ID is missing!");
            return;
          }

          window.location.href =
            "https://accounts.google.com/o/oauth2/v2/auth?" +
            new URLSearchParams({
              client_id: GOOGLE_CLIENT_ID,
              redirect_uri: "http://localhost:3000/google/callback",
              response_type: "id_token",
              scope: "openid email profile",
              nonce: "123xyz",
            }).toString();
        }}
      >
        Sign in with Google
      </button>


      <p className="text-center mt-3">
        Don't have an account? <a href="/register">Register</a>
      </p>
    </div>
  );
}
