import { useState } from "react";
import { register } from "../services/authService";
import { useNavigate } from "react-router-dom";
import FormInput from "../components/FormInput"

export default function Register() {
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [birthDay, setBirthDay] = useState("");
  const [password, setPassword] = useState("");

  const navigate = useNavigate();

  async function handleSubmit(e) {
    e.preventDefault();

    const res = await register({ fullName, email, password, birthDay });

    if (res.success) {
      navigate(`/verify-otp?email=${email}`);
    } else {
      alert(res.message);
    }
  }

  return (
    <div className="container mt-5" style={{ maxWidth: 450 }}>
      <h3 className="text-center mb-4">Register</h3>

      <form onSubmit={handleSubmit}>
        <FormInput label="Full Name" value={fullName} setValue={setFullName} />
        <FormInput label="Email" value={email} setValue={setEmail} />
        <FormInput label="Password" type="password" value={password} setValue={setPassword} />
        
        <FormInput label="Birth Day" type="date" value={birthDay} setValue={setBirthDay} />

        <button className="btn btn-primary w-100">Register</button>
      </form>

      <p className="text-center mt-3">
        Already have an account? <a href="/sign-in">Sign in</a>
      </p>
    </div>
  );
}
