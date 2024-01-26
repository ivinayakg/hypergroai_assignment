import { Button } from "@/components/ui/button";
import { ModeToggle } from "@/components/mode-toggle";
import { useNavigate, Link, useLocation } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import { setInLocalStorage } from "@/utils/localstorage";

function Header() {
  const loginWithGoogleUrl = import.meta.env.VITE_GOOGLE_SIGN_IN;
  const urls = [
    {
      name: "Top",
      url: "/top",
    },
    {
      name: "Home",
      url: "/",
    },
    {
      name: "Favourites",
      url: "/favourites",
    },
  ];
  const { userState } = useMain();
  const navigate = useNavigate();
  const location = useLocation();
  const pathName = location.pathname;

  const logout = () => {
    setInLocalStorage("userToken:hypergroai", null);
    navigate(0);
  };

  return (
    <div className="min-w-full flex justify-between items-center py-5">
      <Link to="/">
        <h1 className="text-3xl">Hypergro AI Assignment</h1>
      </Link>
      <div className="flex justify-center items-center gap-4">
        {userState.login ? (
          <>
            <Button className="border-2 border-primary" variant={"outline"} onClick={logout}>
              Logout
            </Button>
            {urls.map((url, i) => (
              <Link to={url.url} key={i}>
                <Button variant={pathName === url.url ? "default" : "ghost"}>
                  {url.name}
                </Button>
              </Link>
            ))}
          </>
        ) : (
          <Link to={loginWithGoogleUrl}>
            <Button className="gap-2">
              Login With Google <img src="/google.svg" alt="" />
            </Button>
          </Link>
        )}
        <ModeToggle />
      </div>
    </div>
  );
}

export default Header;
