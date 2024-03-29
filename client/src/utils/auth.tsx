import { useNavigate, useSearchParams } from "react-router-dom";
import fetch from "./axios";
import { setInLocalStorage } from "./localstorage";
import { useEffect } from "react";

export default async function (token: string) {
  if (!token) return;
  const res = await fetch.get("/user/self", {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (res.status !== 200) {
    setInLocalStorage("userToken:hypergroai", null);
    return;
  } else {
    setInLocalStorage("userToken:hypergroai", token);
  }
  return res.data;
}

function AuthComponent() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      setInLocalStorage("userToken:hypergroai", token);
    }
    navigate("/");
  });

  return null;
}

export { AuthComponent };
