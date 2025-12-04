import { useSearchParams, useNavigate } from "react-router-dom";
import { verifyOtp } from "../services/authService";
import { useState } from "react";

export default function VerifyOtp() {
  const [params] = useSearchParams();
  const email = params.get("email");

  const [otpCode, setOtpCode] = useState("");
  const navigate = useNavigate();

  async function handleSubmit(e) {
    e.preventDefault();

    const res = await verifyOtp({ email, otpCode });
    if (res.success) {
      alert("Verify success!");
      navigate("/sign-in");
    } else {
      alert(res.message);
    }
  }

  return (
    <div className="container mt-5" style={{ maxWidth: 450 }}>
      <h3 className="text-center mb-4">Verify OTP</h3>

      <form onSubmit={handleSubmit}>
        <div className="mb-3">
          <label>OTP Code</label>
          <input
            className="form-control"
            value={otpCode}
            onChange={(e) => setOtpCode(e.target.value)}
          />
        </div>

        <button className="btn btn-success w-100">Verify</button>
      </form>
    </div>
  );
}
