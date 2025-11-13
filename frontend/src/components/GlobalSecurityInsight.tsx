import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Area, AreaChart } from "recharts";
import { AlertTriangle, TrendingUp, Clock } from "lucide-react";
import { format, subDays } from "date-fns";

type GlobalSecurityStat = {
  date: string;
  count: number;
};

const MOCK_DATA: GlobalSecurityStat[] = Array.from({ length: 30 }, (_, i) => ({
  date: format(subDays(new Date(), 29 - i), "yyyy-MM-dd"),
  count: Math.floor(Math.random() * 50) + 20,
}));

export function GlobalSecurityInsight() {
  const [data, setData] = useState<GlobalSecurityStat[]>(MOCK_DATA);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());

  const fetchCVEData = async () => {
    try {
      setLoading(true);
      const endDate = format(new Date(), "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'");
      const startDate = format(subDays(new Date(), 30), "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'");
      
      const response = await fetch(
        `https://services.nvd.nist.gov/rest/json/cves/2.0?pubStartDate=${startDate}&pubEndDate=${endDate}`
      );

      if (!response.ok) throw new Error("Failed to fetch CVE data");

      const result = await response.json();
      
      // Aggregate CVEs by day
      const cvesByDay = new Map<string, number>();
      result.vulnerabilities?.forEach((vuln: any) => {
        const publishDate = format(new Date(vuln.cve.published), "yyyy-MM-dd");
        cvesByDay.set(publishDate, (cvesByDay.get(publishDate) || 0) + 1);
      });

      const aggregated: GlobalSecurityStat[] = Array.from({ length: 30 }, (_, i) => {
        const date = format(subDays(new Date(), 29 - i), "yyyy-MM-dd");
        return {
          date,
          count: cvesByDay.get(date) || 0,
        };
      });

      setData(aggregated);
      setError(false);
      setLastUpdated(new Date());
    } catch (err) {
      console.error("Failed to fetch CVE data, using mock data:", err);
      setError(true);
      setData(MOCK_DATA);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCVEData();
    
    // Auto-refresh every 6 hours
    const interval = setInterval(fetchCVEData, 6 * 60 * 60 * 1000);
    return () => clearInterval(interval);
  }, []);

  const calculateTrend = () => {
    if (data.length < 14) return 0;
    
    const lastWeek = data.slice(-7).reduce((sum, d) => sum + d.count, 0);
    const previousWeek = data.slice(-14, -7).reduce((sum, d) => sum + d.count, 0);
    
    if (previousWeek === 0) return 0;
    return Math.round(((lastWeek - previousWeek) / previousWeek) * 100);
  };

  const trend = calculateTrend();
  const totalThisWeek = data.slice(-7).reduce((sum, d) => sum + d.count, 0);

  return (
    <Card className="animate-fadeInUp shadow-soft hover:shadow-elegant transition-smooth">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <span className="text-lg">🌐 Global Security Incidents</span>
              {loading && (
                <Clock className="h-4 w-4 text-muted-foreground animate-spin" />
              )}
            </CardTitle>
            <CardDescription className="mt-1">
              Daily reported vulnerabilities and breaches worldwide
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {error && (
          <div className="flex items-center gap-2 p-3 rounded-lg bg-muted/50 border border-border">
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
            <p className="text-xs text-muted-foreground">
              Unable to fetch live data. Displaying sample insights.
            </p>
          </div>
        )}

        <div className="h-[240px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <defs>
                <linearGradient id="colorGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" />
              <XAxis
                dataKey="date"
                tickFormatter={(value) => format(new Date(value), "MMM dd")}
                stroke="hsl(var(--muted-foreground))"
                fontSize={12}
                tickMargin={8}
              />
              <YAxis
                stroke="hsl(var(--muted-foreground))"
                fontSize={12}
                tickMargin={8}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: "hsl(var(--card))",
                  border: "1px solid hsl(var(--border))",
                  borderRadius: "var(--radius)",
                  boxShadow: "var(--shadow-soft)",
                }}
                labelFormatter={(value) => format(new Date(value), "MMMM dd, yyyy")}
                formatter={(value: number) => [value, "Vulnerabilities"]}
              />
              <Area
                type="monotone"
                dataKey="count"
                stroke="hsl(var(--primary))"
                strokeWidth={2}
                fill="url(#colorGradient)"
                animationDuration={800}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>

        <div className="flex items-center justify-between pt-4 border-t border-border">
          <div className="flex items-center gap-2">
            <TrendingUp
              className={`h-4 w-4 ${
                trend > 0 ? "text-destructive" : "text-success"
              }`}
            />
            <span className="text-sm font-medium">
              {trend > 0 ? "↑" : "↓"} {Math.abs(trend)}% {trend > 0 ? "more" : "fewer"} vulnerabilities this week
            </span>
          </div>
          <p className="text-xs text-muted-foreground">
            Last updated: {format(lastUpdated, "HH:mm")}
          </p>
        </div>

        <div className="p-4 rounded-lg bg-secondary/50 border border-border">
          <p className="text-xs text-muted-foreground leading-relaxed">
            <strong className="text-foreground">Security threats grow daily</strong> — D-Vault ensures your credentials stay sovereign and encrypted beyond global attack trends.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
