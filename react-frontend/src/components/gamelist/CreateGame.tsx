"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
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
import axios from "axios";
import { Plus } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

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

export function CreateGameDialog() {
  const [open, setOpen] = useState(false);
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
      return axios.post(
        "http://localhost:8080/games/new",
        {
          name: form.getValues("name"),
          subject: form.getValues("subject"),
          max_player_count: form.getValues("max_player_count"),
          winner_count: form.getValues("winner_count"),
        },
        { withCredentials: true }
      );
    },
    onSuccess: () => {
      setOpen(false);
      form.reset();
      queryClient.invalidateQueries({ queryKey: ["games"] }); // Invalidate the games query to refresh the list
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
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Create New Game
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px] bg-slate-900 text-white border-slate-700">
        {!mutation.isPending ? (
          <>
            <DialogHeader>
              <DialogTitle className="text-xl font-bold">
                Create New Game
              </DialogTitle>

              <DialogDescription className="text-slate-400">
                Set up your game parameters and invite players.
              </DialogDescription>
            </DialogHeader>
            <Form {...form}>
              <form className="space-y-6" onSubmit={handleFormSubmit}>
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-slate-300">
                        Game Name
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Enter game name"
                          {...field}
                          className="bg-slate-800 border-slate-700 text-white"
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
                      <FormLabel className="text-slate-300">
                        Card theme
                      </FormLabel>
                      <FormDescription className="text-slate-400">
                        Choose a theme for your card collection. Our AI will
                        generate a unique deck of cards based on this
                        subject—try options like "Zombie Apocalypse Karaoke" or
                        "Intergalactic Burrito Bash."
                      </FormDescription>
                      <FormControl>
                        <Input
                          placeholder="Enter a card theme"
                          {...field}
                          className="bg-slate-800 border-slate-700 text-white"
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
                      <FormLabel className="text-slate-300">
                        Max Players: {field.value}
                      </FormLabel>
                      <FormDescription className="text-slate-400">
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
                      <FormLabel className="text-slate-300">
                        Winning count: {field.value}
                      </FormLabel>
                      <FormDescription className="text-slate-400">
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

                {/* <FormField
            control={form.control}
            name="is_private"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded-lg border border-slate-700 p-3">
                <div className="space-y-0.5">
                  <FormLabel className="text-slate-300">
                    Private Game
                  </FormLabel>
                  <FormDescription className="text-slate-400">
                    Make this game private with password protection
                  </FormDescription>
                </div>
                <FormControl>
                  <Switch
                    checked={field.value}
                    onCheckedChange={(checked) => {
                      field.onChange(checked);
                      setIsPrivate(checked);
                    }}
                    className="data-[state=checked]:bg-blue-600"
                  />
                </FormControl>
              </FormItem>
            )}
          /> */}
                {/* 
          {isPrivate && (
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-slate-300">Password</FormLabel>
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="Enter password"
                      {...field}
                      className="bg-slate-800 border-slate-700 text-white"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          )} */}

                <DialogFooter>
                  <Button
                    type="submit"
                    className="bg-blue-600 hover:bg-blue-700 text-white"
                  >
                    Create Game
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          </>
        ) : (
          <Spinner />
        )}
      </DialogContent>
    </Dialog>
  );
}
