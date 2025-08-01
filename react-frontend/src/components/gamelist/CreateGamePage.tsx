"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Slider } from "@/components/ui/slider";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/axios";
import { ArrowLeft } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import * as z from "zod";
import { SiteHeader } from "../site-header";

const formSchema = z.object({
  name: z.string().min(3, {
    message: "Game name must be at least 3 characters.",
  }),
  game_type: z.string({
    required_error: "Please select a game type.",
  }),
  subject: z.string().min(3, {
    message: "Collection name must be at least 3 characters.",
  }),
  max_player_count: z.number().min(2).max(6),
  winner_count: z.number().min(3).max(50),
  is_private: z.boolean().default(false).optional(),
  password: z.string().optional(),
});

export function CreateGamePage() {
  const navigate = useNavigate();
  const [isPrivate, setIsPrivate] = useState(false);
  const queryClient = useQueryClient();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      subject: "",
      max_player_count: 3,
      winner_count: 10,
      is_private: false,
      password: "",
    },
  });

  const mutation = useMutation({
    mutationFn: async () => {
      const response = await api.post("/games/new", {
        name: form.getValues("name"),
        subject: form.getValues("subject"),
        max_player_count: form.getValues("max_player_count"),
        winner_count: form.getValues("winner_count"),
      });
      return response.data;
    },
    onSuccess: (data) => {
      form.reset();
      console.log(data);
      navigate(`/games/${data.id}`);
    },
    onError: (error) => {
      console.error(error);
    },
  });

  function onSubmit() {
    mutation.mutate();
  }

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit();
  };

  return (
    <>
      <SiteHeader title="Create New Game" />
      <div className="container mx-auto px-4 py-6">
        <div className="max-w-2xl mx-auto">
          {/* Header */}
          <div className="mb-8">
            <Button
              variant="ghost"
              onClick={() => navigate("/games")}
              className="mb-4"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to games
            </Button>
            <h1 className="text-3xl font-bold mb-2">
              Create New Game
            </h1>
            <p className="text-muted-foreground">
              Set up your game parameters and invite players.
            </p>
          </div>

        {/* Form */}
        {!mutation.isPending ? (
          <Card>
            <CardContent className="p-6">
              <Form {...form}>
                <form className="space-y-6" onSubmit={handleFormSubmit}>
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          Game Name
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder="Enter game name"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="subject"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          Card subject
                        </FormLabel>
                        <FormDescription>
                          Choose a subject for your card collection. Our AI will
                          generate a unique deck of cards based on this
                          subject — try options like "Call of Duty" or
                          "Ancient Mythology."
                        </FormDescription>
                        <FormControl>
                          <Input
                            placeholder="Enter a card theme"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="max_player_count"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          Max Players: {field.value}
                        </FormLabel>
                        <FormDescription>
                          Set the maximum number of players (2-6)
                        </FormDescription>
                        <FormControl>
                          <Slider
                            min={2}
                            max={6}
                            step={1}
                            defaultValue={[field.value]}
                            onValueChange={(vals) => field.onChange(vals[0])}
                            className="py-4"
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="winner_count"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          Winning Count: {field.value}
                        </FormLabel>
                        <FormDescription>
                          Set the amount of points needed to win
                        </FormDescription>
                        <FormControl>
                          <Slider
                            min={3}
                            max={20}
                            step={1}
                            defaultValue={[field.value]}
                            onValueChange={(vals) => field.onChange(vals[0])}
                            className="py-4"
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex gap-4 pt-4">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => navigate("/games")}
                      className="flex-1"
                    >
                      Cancel
                    </Button>
                    <Button
                      type="submit"
                      className="flex-1"
                    >
                      Create Game
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        ) : (
          <div className="flex items-center justify-center py-12">
            <Spinner />
          </div>
        )}
      </div>
    </div>
    </>
  );
} 