import React from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import { Menu } from 'lucide-react';
import { cn } from '@/lib/utils';

interface NavItem {
  href: string;
  label: string;
}

const navItems: NavItem[] = [
  { href: '/', label: 'Home' },
  { href: '/events', label: 'Recent Changes' },
  { href: '/lifecycle', label: 'Resource Lifecycle' },
];

const SidebarContent: React.FC<{ currentPath: string }> = ({ currentPath }) => {
  return (
    <nav className="flex flex-col space-y-1">
      {navItems.map((item) => {
        const isActive = currentPath === item.href;
        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "px-3 py-2 rounded-md text-sm font-medium transition-colors",
              isActive
                ? "bg-primary text-primary-foreground"
                : "text-foreground hover:bg-muted hover:text-foreground"
            )}
          >
            {item.label}
          </Link>
        );
      })}
    </nav>
  );
};

export const Sidebar: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const router = useRouter();

  return (
    <div className="flex min-h-screen">
      <aside className="hidden md:flex md:w-64 md:flex-col md:border-r md:bg-background">
        <div className="flex h-full flex-col px-4 py-6">
          <SidebarContent currentPath={router.pathname} />
        </div>
      </aside>

      <div className="flex flex-1 flex-col">
        <header className="sticky top-0 z-10 flex h-14 items-center gap-4 border-b bg-background px-4 md:hidden">
          <Sheet>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon">
                <Menu className="h-5 w-5" />
                <span className="sr-only">Toggle menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="w-64">
              <div className="mt-4">
                <SidebarContent currentPath={router.pathname} />
              </div>
            </SheetContent>
          </Sheet>
          <h1 className="text-lg font-semibold">Kubernetes Auditing Dashboard</h1>
        </header>

        <main className="flex-1">{children}</main>
      </div>
    </div>
  );
};
