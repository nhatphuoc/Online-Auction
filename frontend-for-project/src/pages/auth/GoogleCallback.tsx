import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { authService } from "../../services/auth";

export default function GoogleCallback() {
    const navigate = useNavigate();
    const ranRef = useRef(false);

    useEffect(() => {
        if (ranRef.current) return;
        ranRef.current = true;

        const hash = window.location.hash.substring(1);
        const params = new URLSearchParams(hash);
        const idToken = params.get("id_token");

        console.log("Received id_token:", idToken);

        if (!idToken) {
            navigate("/login");
            return;
        }

        // Clear hash
        window.history.replaceState(null, "", window.location.pathname);

        authService
            .signInWithGoogle(idToken)
            .then((res) => {
                if (res.success) {
                    navigate("/");
                } else {
                    console.error(res.message);
                    navigate("/login");
                }
            })
            .catch((err) => {
                console.error("Google sign-in error:", err);
                navigate("/login");
            });
    }, []);

    return (
        <div className="container mt-5 text-center">
            <h3>Đang đăng nhập bằng Google...</h3>
        </div>
    );
}