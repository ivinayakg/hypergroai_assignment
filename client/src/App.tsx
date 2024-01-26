import { ThemeProvider } from "@/components/theme-provider";
import "./App.css";
import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom";
import { MainProvider } from "@/components/main-provider";
import { Toaster } from "@/components/ui/toaster";
import { AuthComponent } from "@/utils/auth";
import StocksPage from "@/pages/stocks";
// import CreateShort from "@/pages/CreateShort";
import NotFound from "@/pages/NotFound";
import Header from "@/components/Header";
import FavouriteStocksPage from "./pages/favouriteStocks";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Main />,
    children: [
      {
        path: "",
        element: <StocksPage />,
      },
      {
        path: "/top",
        element: <StocksPage />,
      },
      {
        path: "/favourites",
        element: <FavouriteStocksPage />,
      },
      {
        path: "/not-found/redirect",
        element: <NotFound />,
      },
    ],
  },
  {
    path: "/auth/",
    element: <AuthComponent />,
  },
]);

function Main() {
  return (
    <>
      <Header />
      <Outlet />
    </>
  );
}

function App() {
  return (
    <>
      <MainProvider>
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
          <RouterProvider router={router} />
          <Toaster />
        </ThemeProvider>
      </MainProvider>
    </>
  );
}

export default App;
