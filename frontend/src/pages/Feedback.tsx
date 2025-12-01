import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { toast } from "@/hooks/use-toast";
import { Send, Loader2, Shield } from "lucide-react";
import { DashboardLayout } from "@/components/DashboardLayout";

// Stub function for sending feedback
async function sendFeedback1(payload: any) {
  console.log("Feedback payload:", payload);
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 1000));
}

export function Feedback1() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [feedbackType, setFeedbackType] = useState<string>("");
  const [message, setMessage] = useState("");
  const [includeLogs, setIncludeLogs] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const maxMessageLength = 1000;
  const isMessageEmpty = message.trim().length === 0;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isMessageEmpty) {
      toast({
        title: "Message required",
        description: "Please enter your feedback message.",
        variant: "destructive",
      });
      return;
    }

    setIsSubmitting(true);

    try {
      const payload = {
        name: name || null,
        email: email || null,
        feedbackType,
        message,
        includeLogs,
        timestamp: new Date().toISOString(),
      };

      await sendFeedback(payload);

      toast({
        title: "Feedback sent!",
        description: "Thank you for your feedback. We appreciate it!",
      });

      // Reset form
      setName("");
      setEmail("");
      setFeedbackType("");
      setMessage("");
      setIncludeLogs(false);
    } catch (error) {
      toast({
        title: "Failed to send feedback",
        description: "Something went wrong. Please try again later.",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
        <DashboardLayout>
    <div className="min-h-screen bg-background flex items-center justify-center p-6">
      <Card className="w-full max-w-2xl shadow-soft">
        <CardHeader>
          <CardTitle className="text-2xl">Send Feedback</CardTitle>
          <CardDescription>
            Help us improve <strong><small>ANKHORA</small></strong> by sharing your thoughts, reporting bugs, or suggesting features.
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Name */}
            <div className="space-y-2">
              <Label htmlFor="name">Name (optional)</Label>
              <Input
                id="name"
                type="text"
                placeholder="Your name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={isSubmitting}
              />
            </div>

            {/* Email */}
            <div className="space-y-2">
              <Label htmlFor="email">Email (optional)</Label>
              <Input
                id="email"
                type="email"
                placeholder="your.email@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={isSubmitting}
              />
            </div>

            {/* Feedback Type */}
            <div className="space-y-2">
              <Label htmlFor="feedbackType">Feedback Type</Label>
              <Select value={feedbackType} onValueChange={setFeedbackType} disabled={isSubmitting}>
                <SelectTrigger id="feedbackType">
                  <SelectValue placeholder="Select a type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bug">Bug</SelectItem>
                  <SelectItem value="feature">Feature Request</SelectItem>
                  <SelectItem value="improvement">Improvement</SelectItem>
                  <SelectItem value="general">General</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Message */}
            <div className="space-y-2">
              <Label htmlFor="message">
                Message <span className="text-destructive">*</span>
              </Label>
              <Textarea
                id="message"
                placeholder="Tell us what's on your mind..."
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                disabled={isSubmitting}
                maxLength={maxMessageLength}
                rows={6}
                className="resize-none"
              />
              <div className="flex justify-between items-center text-xs text-muted-foreground">
                <span>Required field</span>
                <span>
                  {message.length} / {maxMessageLength}
                </span>
              </div>
            </div>

            <Separator />

            {/* Include System Logs */}
            <div className="flex items-center justify-between space-x-2">
              <div className="space-y-0.5">
                <Label htmlFor="includeLogs" className="cursor-pointer">
                  Include system logs
                </Label>
                <p className="text-sm text-muted-foreground">
                  Help us diagnose issues by including technical logs
                </p>
              </div>
              <Switch
                id="includeLogs"
                checked={includeLogs}
                onCheckedChange={setIncludeLogs}
                disabled={isSubmitting}
              />
            </div>

            <Separator />

            {/* Submit Button */}
            <Button
              type="submit"
              className="w-full"
              disabled={isMessageEmpty || isSubmitting}
            >
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Sending...
                </>
              ) : (
                <>
                  <Send className="mr-2 h-4 w-4" />
                  Send Feedback
                </>
              )}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>

        </DashboardLayout>
  );
}

async function sendFeedback(payload: any) {
  console.log("Feedback payload:", payload);
  await new Promise(resolve => setTimeout(resolve, 1000));
}

export default function Feedback() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [feedbackType, setFeedbackType] = useState<string>("");
  const [message, setMessage] = useState("");
  const [includeLogs, setIncludeLogs] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const maxMessageLength = 1000;
  const isMessageEmpty = message.trim().length === 0;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isMessageEmpty) {
      toast({
        title: "Message required",
        description: "Please enter your feedback message.",
        variant: "destructive",
      });
      return;
    }

    setIsSubmitting(true);

    try {
      const payload = {
        name: name || null,
        email: email || null,
        feedbackType,
        message,
        includeLogs,
        timestamp: new Date().toISOString(),
      };

      await sendFeedback(payload);

      toast({
        title: "Feedback sent!",
        description: "Thank you for helping us improve Ankhora.",
      });

      setName("");
      setEmail("");
      setFeedbackType("");
      setMessage("");
      setIncludeLogs(false);
    } catch (error) {
      toast({
        title: "Failed to send feedback",
        description: "Something went wrong. Please try again later.",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <DashboardLayout>
      <div className="min-h-screen flex items-center justify-center p-8 bg-gradient-to-b from-zinc-50 to-zinc-100 dark:from-zinc-950 dark:to-zinc-900">
        <Card className="w-full max-w-2xl border-none shadow-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-900/50">
          <CardHeader className="text-center space-y-6 pb-2">
            <div className="mx-auto w-16 h-16 rounded-2xl bg-gradient-to-br from-primary/20 to-amber-500/20 backdrop-blur-sm shadow-lg flex items-center justify-center">
              <Shield className="h-8 w-8 text-primary" />
            </div>
            <div className="space-y-2">
              <CardTitle className="text-3xl font-semibold tracking-tight bg-gradient-to-r from-foreground to-primary/80 bg-clip-text text-transparent">
                Send Feedback
              </CardTitle>
              <CardDescription className="text-lg text-muted-foreground/80 leading-relaxed max-w-md mx-auto">
                Help us improve <span className="font-semibold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">ANKHORA</span>
              </CardDescription>
            </div>
          </CardHeader>

          <CardContent className="p-8">
            <form onSubmit={handleSubmit} className="space-y-7">
              {/* Optional Fields Row */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2 group">
                  <Label htmlFor="name" className="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors">
                    Name (optional)
                  </Label>
                  <Input
                    id="name"
                    type="text"
                    className="h-12 rounded-xl border-zinc-200 dark:border-zinc-700 focus:ring-2 focus:ring-primary/30 transition-all backdrop-blur-sm"
                    placeholder="Your name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>

                <div className="space-y-2 group">
                  <Label htmlFor="email" className="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors">
                    Email (optional)
                  </Label>
                  <Input
                    id="email"
                    type="email"
                    className="h-12 rounded-xl border-zinc-200 dark:border-zinc-700 focus:ring-2 focus:ring-primary/30 transition-all backdrop-blur-sm"
                    placeholder="your.email@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>
              </div>

              {/* Feedback Type */}
              <div className="space-y-2">
                <Label htmlFor="feedbackType" className="text-sm font-medium">
                  Feedback Type
                </Label>
                <Select value={feedbackType} onValueChange={setFeedbackType} disabled={isSubmitting}>
                  <SelectTrigger className="h-12 rounded-xl border-zinc-200 dark:border-zinc-700 focus:ring-2 focus:ring-primary/30 data-[state=open]:ring-2 data-[state=open]:ring-primary/30 backdrop-blur-sm">
                    <SelectValue placeholder="Select a type" />
                  </SelectTrigger>
                  <SelectContent className="backdrop-blur-sm bg-white/80 dark:bg-zinc-900/80 border-zinc-200/50 dark:border-zinc-700/50">
                    <SelectItem value="bug" className="hover:bg-primary/5 focus:bg-primary/5">Bug</SelectItem>
                    <SelectItem value="feature" className="hover:bg-primary/5 focus:bg-primary/5">Feature Request</SelectItem>
                    <SelectItem value="improvement" className="hover:bg-primary/5 focus:bg-primary/5">Improvement</SelectItem>
                    <SelectItem value="general" className="hover:bg-primary/5 focus:bg-primary/5">General</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Message */}
              <div className="space-y-3">
                <Label htmlFor="message" className="text-sm font-medium flex items-center gap-1">
                  Message <span className="text-destructive text-xs">*</span>
                </Label>
                <Textarea
                  id="message"
                  className="h-32 resize-none rounded-2xl border-zinc-200 dark:border-zinc-700 focus:ring-2 focus:ring-primary/30 p-4 transition-all backdrop-blur-sm min-h-[120px]"
                  placeholder="Tell us what's on your mind..."
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  disabled={isSubmitting}
                  maxLength={maxMessageLength}
                />
                <div className="flex justify-between items-center text-xs text-muted-foreground">
                  <span>Required field</span>
                  <span className="font-mono font-medium">
                    {message.length} / {maxMessageLength}
                  </span>
                </div>
              </div>

              <Separator className="border-zinc-200/50 dark:border-zinc-700/50 my-1" />

              {/* Include Logs */}
              <div className="group flex items-center justify-between p-4 rounded-2xl bg-white/30 dark:bg-zinc-800/30 backdrop-blur-sm border border-zinc-200/30 hover:bg-white/50 dark:hover:bg-zinc-800/50 transition-all">
                <div className="space-y-1">
                  <Label htmlFor="includeLogs" className="text-sm font-medium text-foreground cursor-pointer">
                    Include system logs
                  </Label>
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    Help us diagnose issues faster with technical details
                  </p>
                </div>
                <Switch
                  id="includeLogs"
                  checked={includeLogs}
                  onCheckedChange={setIncludeLogs}
                  disabled={isSubmitting}
                  className="data-[state=checked]:bg-primary"
                />
              </div>

              <Separator className="border-zinc-200/50 dark:border-zinc-700/50 my-1" />

              {/* Submit */}
              <Button
                type="submit"
                className="w-full h-14 rounded-2xl text-lg font-semibold bg-gradient-to-r from-primary to-amber-500 hover:from-primary/90 hover:to-amber-500/90 shadow-xl hover:shadow-primary/20 hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none"
                disabled={isMessageEmpty || isSubmitting}
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="mr-3 h-5 w-5 animate-spin" />
                    Sending feedback...
                  </>
                ) : (
                  <>
                    <Send className="mr-3 h-5 w-5" />
                    Send Feedback
                  </>
                )}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  );
}
