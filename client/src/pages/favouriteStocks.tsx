import { useLocation, useNavigate } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import { useEffect, useState } from "react";
import fetch from "@/utils/axios";
import { useToast } from "@/components/ui/use-toast";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

function FavouriteStocksPage() {
  const { userState } = useMain();
  const [userFavourites, setUserFavourites] = useState([]);

  const { toast } = useToast();
  const Headings = [
    "Serial",
    "Date",
    "Stock Name",
    "Gain",
    "Open",
    "High",
    "Low",
    "Close",
    "Last",
    "Action"
  ];
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      const res = await fetch.get("/favourite", {
        headers: { Authorization: `Bearer ${userState.token}` },
      });
      const data = res.data;
      setUserFavourites(data.data);
    })();
    if (!userState.login) navigate("/");
  }, [userState.login, location.pathname]);

  const unfavouriteStockHandler = async (stockCode: string) => {
    const res = await fetch.delete(`/favourite/${stockCode}`);
    if (res.status >= 400) {
      toast({
        title: "Something went wrong",
        duration: 2000,
      });
      return;
    }
    setUserFavourites((prev: any) => {
      return prev.filter((stock: any) => stock.stockCode !== stockCode);
    });
  };

  return (
    <div className="py-5">
      <Table>
        <TableHeader>
          <TableRow>
            {Headings.map((v, i) => {
              return (
                <TableHead className="w-[100px]" key={i}>
                  {v}
                </TableHead>
              );
            })}
          </TableRow>
        </TableHeader>
        <TableBody>
          {userFavourites.length ? (
            userFavourites.map((stock: any, i) => {
              return (
                <TableRow key={stock._id}>
                  <TableCell className="font-medium text-left">
                    {i + 1}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {stock.date}
                  </TableCell>
                  <TableCell className="font-medium text-left hover:cursor-pointer">
                    {stock.stockName}
                    <p className="text-xs opacity-80">/{stock.stockCode}</p>
                  </TableCell>
                  <TableCell
                    className={`font-medium text-left ${
                      stock.gain > 0 ? "text-green-400" : "text-red-400"
                    }`}
                  >
                    {stock.gain}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {stock.open}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {stock.high}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {stock.low}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {stock.close}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {stock.last}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    <Button
                      variant={"outline"}
                      onClick={() => unfavouriteStockHandler}
                    >
                      Unfavourite
                    </Button>
                  </TableCell>
                </TableRow>
              );
            })
          ) : (
            <></>
          )}
        </TableBody>
      </Table>
    </div>
  );
}

export default FavouriteStocksPage;
