import { IconDashboard } from "@tabler/icons-react";
import { Link } from "react-router-dom";

import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { CirclePlus } from "lucide-react";

export function NavMain() {
  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Dashboard">
              <Link to="/games">
                <IconDashboard />
                <span>Games</span>
              </Link>
            </SidebarMenuButton>
            <SidebarMenuButton asChild tooltip="Dashboard">
              <Link to="/games/new">
                <CirclePlus />
                <span>Create a game</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
}
