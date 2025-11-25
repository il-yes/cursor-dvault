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
import { Send, Loader2 } from "lucide-react";
import { DashboardLayout } from "@/components/DashboardLayout";

// Stub function for sending feedback
async function sendFeedback(payload: any) {
  console.log("Feedback payload:", payload);
  // Simulate API delay
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
            Help us improve D-Vault by sharing your thoughts, reporting bugs, or suggesting features.
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
