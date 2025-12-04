import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

export default function GoogleCallback() {
  const navigate = useNavigate();

  useEffect(() => {
    // Google trả id_token trong hash (#), không phải query
    const hashParams = new URLSearchParams(window.location.hash.substring(1));
    const idToken = hashParams.get("id_token");

    if (!idToken) {
        console.error("Không tìm thấy id_token!");
        navigate("/sign-in");
        return;
    }

    axios.post("http://localhost:8080/auth/sign-in/google", { idToken })
        .then(res => {
        localStorage.setItem("accessToken", res.data.accessToken);
        localStorage.setItem("refreshToken", res.data.refreshToken);
        navigate("/");
        })
        .catch(err => {
        console.log(err);
        navigate("/sign-in");
        });
    }, [navigate]);


  return (
    <div className="container mt-5 text-center">
      <h3>Đang đăng nhập bằng Google...</h3>
    </div>
  );
}
