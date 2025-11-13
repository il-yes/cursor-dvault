import { CheckCircle2, Clock, XCircle } from "lucide-react";
import { Badge } from "@/components/ui/badge";

type StatusType = "verified" | "pending" | "failed";

interface StatusBadgeProps {
  status: StatusType;
  label: string;
}

export function StatusBadge({ status, label }: StatusBadgeProps) {
  const config = {
    verified: {
      icon: CheckCircle2,
      className: "bg-success/10 text-success border-success/30 hover:bg-success/20",
    },
    pending: {
      icon: Clock,
      className: "bg-warning/10 text-warning border-warning/30 hover:bg-warning/20",
    },
    failed: {
      icon: XCircle,
      className: "bg-destructive/10 text-destructive border-destructive/30 hover:bg-destructive/20",
    },
  };

  const { icon: Icon, className } = config[status];

  return (
    <Badge variant="outline" className={`${className} gap-1.5 transition-all`}>
      <Icon className="h-3 w-3" />
      {label}
    </Badge>
  );
}
