import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { authService } from "../../services/auth";
import { useAuthStore } from '../../stores/auth.store';

export default function GoogleCallback() {
    const navigate = useNavigate();
    const ranRef = useRef(false);
    const setUser = useAuthStore((state) => state.setUser);

    useEffect(() => {
        if (ranRef.current) return;
        ranRef.current = true;

        const hash = window.location.hash.substring(1);
        const params = new URLSearchParams(hash);
        const idToken = params.get("id_token");

        if (!idToken) {
            navigate("/login");
            return;
        }

        // Clear hash
        window.history.replaceState(null, "", window.location.pathname);

        authService
            .signInWithGoogle(idToken)
            .then(async (res) => {
                if (res.success) {
                    const user = await authService.getCurrentUser();
                    setUser(user);
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