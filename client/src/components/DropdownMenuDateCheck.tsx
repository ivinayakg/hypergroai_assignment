import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";

export function DropdownMenuDateCheck({ date, setDate }: { date: any; setDate: any; }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline">Date</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56">
        <DropdownMenuLabel>Select Date</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {date?.options?.map((option: any) => {
          return (
            <DropdownMenuCheckboxItem
              key={option._id}
              checked={option.dataDate === date.current}
              onCheckedChange={() => {
                setDate((prev: any) => ({ ...prev, current: option.dataDate }));
              }}
            >
              {option.dataDate}
            </DropdownMenuCheckboxItem>
          );
        })}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
