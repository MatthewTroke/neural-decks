import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Link } from "react-router";

export function SiteHeader({ title }: { title: string }) {
  return (
    <header className="flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height) z-10">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1" />
        <Separator
          orientation="vertical"
          className="mx-2 data-[orientation=vertical]:h-4"
        />
        <h1 className="text-base font-medium">{title}</h1>
      </div>
      <div className="block md:hidden w-full flex justify-end mx-2">
        <Link to="/" className="flex items-center gap-2">
          <span className="text-md font-bold">Neural Decks</span>
          <div className="inline-block rounded-lg bg-purple-900/20 px-3 py-1 text-sm text-purple-400">
            Beta
          </div>
        </Link>
      </div>
    </header>
  );
}
