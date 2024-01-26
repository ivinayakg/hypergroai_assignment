import { useLocation, useNavigate } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import { useEffect, useRef, useState } from "react";
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
import { Input } from "@/components/ui/input";
import { DropdownMenuDateCheck } from "../components/DropdownMenuDateCheck";
import { Checkbox } from "@/components/ui/checkbox";

function StocksPage() {
  const { userState } = useMain();
  const [stockInfoData, setStockInfoData] = useState([]);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [date, setDate] = useState({ current: "", options: [] });
  const searchDebounce = useRef<any>();
  const [userFavouritesCodes, setUserFavouritesCodes] = useState(new Set());

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
    "Favourite",
  ];
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      const res = await fetch.get("/favourite/codes", {
        headers: { Authorization: `Bearer ${userState.token}` },
      });
      const data = res.data;
      setUserFavouritesCodes(new Set(data.data));
    })();
  }, [userState.login]);

  useEffect(() => {
    (async () => {
      const res = await fetch.get("/migration", {
        headers: { Authorization: `Bearer ${userState.token}` },
      });
      const data = res.data;
      setDate((prev) => ({
        ...prev,
        options: data.data,
        current: data.data[0].dataDate,
      }));
    })();
  }, [userState.login]);

  useEffect(() => {
    let url = `/stock?page=${page}&date=${date.current}&s=${search}`;
    if (!userState.login) {
      url = `/stock/unverified`;
    } else if (location.pathname == "/top") {
      url = `/stock/top?page=${page}&date=${date.current}&s=${search}`;
    }

    (async () => {
      const res = await fetch.get(url, {
        headers: { Authorization: `Bearer ${userState.token}` },
      });
      const data = res.data;
      setStockInfoData(data.data);
      if (data.length === 0) {
        toast({
          title: "Something went wrong",
          duration: 2000,
        });
      }
    })();

    if (!userState.login) navigate("/");
  }, [page, date.current, search, userState.login, location.pathname]);

  const changeSearch = (e: any) => {
    let value = e.target.value;
    clearTimeout(searchDebounce.current);
    let i = setTimeout(() => {
      setSearch(value);
      setPage(1);
    }, 1500);
    searchDebounce.current = i;
  };

  const favouriteStockHandler = async (stockCode: string) => {
    const res = await fetch.post(
      "/favourite",
      { stockCode },
      {
        headers: { Authorization: `Bearer ${userState.token}` },
      }
    );
    const data = res.data;

    setUserFavouritesCodes(new Set(data.data.stocks));
  };

  return (
    <div className="py-5 flex flex-col justify-center items-center gap-5">
      {userState.login && (
        <div className="flex justify-center items-center gap-4">
          <DropdownMenuDateCheck date={date} setDate={setDate} />
          <Input
            placeholder="Stock Name"
            type="text"
            className="max-w-64"
            onChange={changeSearch}
          />
          <Button
            onClick={() => setPage((prev: any) => prev - 1)}
            disabled={page === 1}
          >
            Prev
          </Button>
          <Button
            onClick={() => setPage((prev: any) => prev + 1)}
            disabled={stockInfoData.length === 0}
          >
            Next
          </Button>
        </div>
      )}

      <Table>
        <TableHeader>
          <TableRow>
            {Headings.map((v, i) => {
              if (!userState.login && v === "Favourite") return;
              return (
                <TableHead className="w-[100px]" key={i}>
                  {v}
                </TableHead>
              );
            })}
          </TableRow>
        </TableHeader>
        <TableBody>
          {stockInfoData.length ? (
            stockInfoData.map((stock: any, i) => {
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
                  {userState.login && (
                    <TableCell className="font-medium text-left">
                      {userFavouritesCodes.has(stock.stockCode) ? (
                        <Checkbox
                          id={stock.stockCode}
                          checked={true}
                          disabled
                        />
                      ) : (
                        <Button
                          variant={"default"}
                          onClick={() => favouriteStockHandler(stock.stockCode)}
                        >
                          Favourite
                        </Button>
                      )}
                    </TableCell>
                  )}
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

export default StocksPage;
