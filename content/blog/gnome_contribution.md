---
title: "Contributing to Gnome"
date: "2026-07-28"
author: "Oliver"
tags: ["linux", "kernel", "dsl"]
---

Hello! Welcome to the last post of DSL26!

In this final phase, our goal was to select a Linux-related open source project (such as KDE, GNOME, Git, Kworflow, etc) and 
work on an issue/collection of issues suited to Newcomers. Since the beginning, me and my colleague were eager to contribute
to the GNOME project. Personally, it's what i use daily :)

Looking through their GitLab instance, we've found some different projects and issues tagged for "Newcomers", and we (among some 
other students, coincidentally), decided for an issue related to the **GNOME** Calendar application.

## The problem 

The issue, which you can see more of [here](https://gitlab.gnome.org/GNOME/gnome-calendar/-/work_items/1092), at a first glance, was that changing events in the month view (a simple dashboard) from one month to another does not enable the whole month grid cell to act as a drop target. Meaning: if you have a lot of events in one month and you wanted to move one more event, you'd wrestle a few seconds before finding out how to drag the event correctly to that cell. 

Pretty weird, right? The following GIF, produced by my colleague, may clear up what is happening.

Link: https://kaikycintra.com/static/gnome_calendar_broken.mp4

## All aboard

Our next step was to find out how could we find the code related to the issue, and work on the project. The **Welcome to GNOME** page gives us some directions to start out and their Code of Conduct. We can also find information specific to the Calendar app [here](https://welcome.gnome.org/en/app/Calendar/). We've also researched some info regarding general [GNOME development](https://developer.gnome.org/documentation/index.html), since they have shared design patterns across applications, defined by their Human Interface Guidelines. 

Initially, it is recommended to use the **GNOME Builder** application. It serves as an IDE for GNOME-related projects, with support for essential technologies such as GTK+, GLib, etc. At a first glance, it would seem that working purely with Nvim would be more straightforward and simple, but believe me, automating the build stage was way better than I thought. Most issues I had with the App is that I had to tackle app dependencies manually through flatpak, since it needs some nightly versions that, for some weird reason, were not download on the first popup message I received.

Working through the issue, we've found the *Drag and Drop* page on Development [Docs](https://developer.gnome.org/documentation/tutorials/drag-and-drop.html), which has some snippets useful for laying out a foundation, or mental model. 

If you go to the *week view** instead of the *month view**, you can see that the issue does not exist there. So, our plan was to compare week's code implementation with month's implementation as an inspiration to our work. We then found the files that tackled drag and drop on the month, such as `src/gui/views/gcal-month-cell.c` and, for comparison, `src/gui/views/gcal-week-grid.c`. Below is the code for month's implementation as an example, for the drop approach (you can skip most of it, it's mostly for reference)

```C 
// gcal-month-cell.c
static void
move_event (GcalMonthCell         *self,
            GcalEvent             *event,
            GcalRecurrenceModType  mod_type)
{

  g_autoptr (GcalEvent) changed_event = NULL;
  g_autoptr (GDateTime) start_dt = NULL;
  GTimeSpan timespan = 0;
  GcalContext *context;
  GDateTime *end_dt;
  gint diff;
  gint start_month, current_month;
  gint start_year, current_year;

  GCAL_ENTRY;
 */

static inline gint
get_parent_bin_height (GtkWidget *widget)
{
  GtkWidget *parent;

  parent = gtk_widget_get_parent (widget);
  g_assert (ADW_IS_BIN (parent));

  return gtk_widget_get_height (parent);
}

static void
update_style_flags (GcalMonthCell *self)
{
  g_autoptr (GDateTime) today = NULL;
  gint weekday;

  /* Today */
  today = g_date_time_new_now_local ();

  if (gcal_date_time_compare_date (self->date, today) == 0)
    gtk_widget_add_css_class (GTK_WIDGET (self), "today");
  else
    gtk_widget_remove_css_class (GTK_WIDGET (self), "today");

  weekday = g_date_time_get_day_of_week (self->date);
  if (is_workday (weekday))
    gtk_widget_add_css_class (GTK_WIDGET (self), "workday");
  else
    gtk_widget_remove_css_class (GTK_WIDGET (self), "workday");
}

static void
add_month_separators (GcalMonthCell *self)
{
  gint day_of_month;

  gtk_widget_remove_css_class (GTK_WIDGET (self), "separator-top");
  gtk_widget_remove_css_class (GTK_WIDGET (self), "separator-side");

  day_of_month = g_date_time_get_day_of_month (self->date);
  if (day_of_month > 1 && day_of_month <= N_WEEKDAYS)
    {
      gtk_widget_add_css_class (GTK_WIDGET (self), "separator-top");
    }
  else if (day_of_month == 1)
    {
      gtk_widget_add_css_class (GTK_WIDGET (self), "separator-side");
      gtk_widget_add_css_class (GTK_WIDGET (self), "separator-top");
    }
}

static void
move_event (GcalMonthCell         *self,
            GcalEvent             *event,
            GcalRecurrenceModType  mod_type)
{

  g_autoptr (GcalEvent) changed_event = NULL;
  g_autoptr (GDateTime) start_dt = NULL;
  GTimeSpan timespan = 0;
  GcalContext *context;
  GDateTime *end_dt;
  gint diff;
  gint start_month, current_month;
  gint start_year, current_year;

  GCAL_ENTRY;

  changed_event = gcal_event_new_from_event (event);

  /* Move the event's date */
  start_dt = gcal_event_get_date_start (changed_event);
  end_dt = gcal_event_get_date_end (changed_event);

  start_month = g_date_time_get_month (start_dt);
  start_year = g_date_time_get_year (start_dt);

  current_month = g_date_time_get_month (self->date);
  current_year = g_date_time_get_year (self->date);

  timespan = g_date_time_difference (end_dt, start_dt);

  start_dt = g_date_time_add_full (start_dt,
                                   current_year - start_year,
                                   current_month - start_month,
                                   0, 0, 0, 0);

  diff = gcal_date_time_compare_date (self->date, start_dt);

  if (diff != 0)
    {
      g_autoptr (GDateTime) new_start = g_date_time_add_days (start_dt, diff);

      gcal_event_set_date_start (changed_event, new_start);

      /* The event may have a NULL end date, so we have to check it here */
      if (end_dt)
        {
          g_autoptr (GDateTime) new_end = g_date_time_add (new_start, timespan);

          gcal_event_set_date_end (changed_event, new_end);
        }

      context = gcal_application_get_context (GCAL_DEFAULT_APPLICATION);
      gcal_manager_update_event (gcal_context_get_manager (context), changed_event, mod_type);
    }
}

<code redacted here ...>

static gboolean
on_drop_target_accept_cb (GtkDropTarget *drop_target,
                          GdkDrop       *drop,
                          GcalMonthCell *self)
{
  GCAL_ENTRY;

  if ((gdk_drop_get_actions (drop) & gtk_drop_target_get_actions (drop_target)) == 0)
    GCAL_RETURN (FALSE);

  if (!gdk_content_formats_contain_gtype (gdk_drop_get_formats (drop), GCAL_TYPE_EVENT_WIDGET))
    GCAL_RETURN (FALSE);

  GCAL_RETURN (TRUE);
}

static void
on_ask_recurrence_response_cb (GcalEvent             *event,
                               GcalRecurrenceModType  mod_type,
                               gpointer               user_data)
{
  GcalMonthCell *self = GCAL_MONTH_CELL (user_data);
      const gchar *temp_str;  /* unowned */

      icon_name = gcal_weather_info_get_icon_name (weather_info);
      temp_str = gcal_weather_info_get_temperature (weather_info);

      gtk_image_set_from_icon_name (self->weather_icon, icon_name);
      gtk_label_set_text (self->temp_label, temp_str);

      /* Use a short month name label to avoid conflicting with the weather forecast's labels */
      if (day_of_month == 1)
        {
          g_autofree gchar *month_name = g_date_time_format (self->date, "%b");
          gtk_label_set_text (self->day_label, month_name);
        }
    }
  else
    {
      gtk_image_clear (self->weather_icon);
      gtk_label_set_text (self->temp_label, "");

      /* No risk of conflicting with the weather forecast labels in their absence, use the full month name label */
      if (day_of_month == 1)
        {
          g_autofree gchar *month_name = g_date_time_format (self->date, "%OB");
          gtk_label_set_text (self->day_label, month_name);
        }
    }
}


/*
 * Callbacks
 */

static void
day_changed_cb (GcalClock     *clock,
                GcalMonthCell *self)
{
  update_style_flags (self);
}

static void
overflow_button_clicked_cb (GtkWidget     *button,
                            GcalMonthCell *self)
{
  g_signal_emit (self, signals[SHOW_OVERFLOW], 0, button);
}

static gboolean
on_drop_target_accept_cb (GtkDropTarget *drop_target,
                          GdkDrop       *drop,
                          GcalMonthCell *self)
{
  GCAL_ENTRY;

  if ((gdk_drop_get_actions (drop) & gtk_drop_target_get_actions (drop_target)) == 0)
    GCAL_RETURN (FALSE);

  if (!gdk_content_formats_contain_gtype (gdk_drop_get_formats (drop), GCAL_TYPE_EVENT_WIDGET))
    GCAL_RETURN (FALSE);

  GCAL_RETURN (TRUE);
}

static void
on_ask_recurrence_response_cb (GcalEvent             *event,
                               GcalRecurrenceModType  mod_type,
                               gpointer               user_data)
{
  GcalMonthCell *self = GCAL_MONTH_CELL (user_data);

  if (mod_type != GCAL_RECURRENCE_MOD_NONE)
    move_event (self, event, mod_type);
}

static gboolean
on_drop_target_drop_cb (GtkDropTarget *drop_target,
                        const GValue  *value,
                        gdouble        x,
                        gdouble        y,
                        GcalMonthCell *self)
{
  GcalEventWidget *event_widget;
  GcalEvent *event;

  GCAL_ENTRY;

  if (!G_VALUE_HOLDS (value, GCAL_TYPE_EVENT_WIDGET))
    GCAL_RETURN (FALSE);

  event_widget = g_value_get_object (value);

  event = gcal_event_widget_get_event (event_widget);

  if (gcal_event_has_recurrence (event))
    {
      gcal_utils_ask_recurrence_modification_type (GTK_WIDGET (self),
                                                   event,
                                                   FALSE,
                                                   on_ask_recurrence_response_cb,
                                                   self);
    }
  else
    {
      move_event (self, event, GCAL_RECURRENCE_MOD_THIS_ONLY);
    }

  GCAL_RETURN (TRUE);
}

```

and below, for week implementation:

```C
static void
move_event_to_cell (GcalWeekGrid          *self,
                    GcalEvent             *event,
                    guint                  cell,
                    GcalRecurrenceModType  mod_type)
{
  GcalContext *context = gcal_application_get_context (GCAL_DEFAULT_APPLICATION);
  g_autoptr (GDateTime) week_start = NULL;
  g_autoptr (GDateTime) dnd_date = NULL;
  g_autoptr (GDateTime) new_end = NULL;
  g_autoptr (GcalEvent) changed_event = NULL;
  GTimeSpan timespan = 0;

  /* RTL languages swap the drop cell column */
  if (gtk_widget_get_direction (GTK_WIDGET (self)) == GTK_TEXT_DIR_RTL)
    {
      gint column, row;

      column = cell / (MINUTES_PER_DAY / 30);
      row = cell - column * 48;

      cell = (N_WEEKDAYS - 1 - column) * 48 + row;
    }
      start = gcal_date_time_add_floating_minutes (week_start, start_cell * 30);
      end = gcal_date_time_add_floating_minutes (week_start, (end_cell + 1) * 30);
    }
  else
    {
      guint rtl_start_cell, rtl_end_cell, rtl_column;

      /* Fix the minute */
      rtl_column = N_WEEKDAYS - 1 - column;
      rtl_start_cell = start_cell + (rtl_column - column) * 48;
      rtl_end_cell = (rtl_column * MINUTES_PER_DAY + minute) / 30;

      start = gcal_date_time_add_floating_minutes (week_start, rtl_start_cell * 30);
      end = gcal_date_time_add_floating_minutes (week_start, (rtl_end_cell + 1) * 30);
    }

  local_x = round ((column + 0.5) * (gtk_widget_get_width (GTK_WIDGET (self)) / (float) N_WEEKDAYS));
  local_y = (minute + 15) * minute_height;

  if (!gtk_widget_compute_point (GTK_WIDGET (self),
                                 weekview,
                                 &GRAPHENE_POINT_INIT (local_x, local_y),
                                 &out))
    g_assert_not_reached ();

  range = gcal_range_new (start, end, GCAL_RANGE_DEFAULT);
  gcal_view_create_event (GCAL_VIEW (weekview), range, out.x, out.y);

  gtk_event_controller_set_propagation_phase (self->motion_controller, GTK_PHASE_NONE);
}

static gint
get_dnd_cell (GcalWeekGrid *self,
              gdouble       x,
              gdouble       y)
{
  gdouble column_width, cell_height;
  gint column, row;

  column_width = gtk_widget_get_width (GTK_WIDGET (self)) / (float) N_WEEKDAYS;
  cell_height = gtk_widget_get_height (GTK_WIDGET (self)) / 48.0;
  column = floor (x / column_width);
  row = y / cell_height;

  return column * 48 + row;
}

static void
move_event_to_cell (GcalWeekGrid          *self,
                    GcalEvent             *event,
                    guint                  cell,
                    GcalRecurrenceModType  mod_type)
{
  GcalContext *context = gcal_application_get_context (GCAL_DEFAULT_APPLICATION);
  g_autoptr (GDateTime) week_start = NULL;
  g_autoptr (GDateTime) dnd_date = NULL;
  g_autoptr (GDateTime) new_end = NULL;
  g_autoptr (GcalEvent) changed_event = NULL;
  GTimeSpan timespan = 0;

  /* RTL languages swap the drop cell column */
  if (gtk_widget_get_direction (GTK_WIDGET (self)) == GTK_TEXT_DIR_RTL)
    {
      gint column, row;

      column = cell / (MINUTES_PER_DAY / 30);
      row = cell - column * 48;

      cell = (N_WEEKDAYS - 1 - column) * 48 + row;
    }

  changed_event = gcal_event_new_from_event (event);
  week_start = gcal_date_time_get_start_of_week (self->active_date);
  dnd_date = gcal_date_time_add_floating_minutes (week_start, cell * 30);

  /*
   * Calculate the diff between the dropped cell and the event's start date,
   * so we can update the end date accordingly.
   */
  timespan = g_date_time_difference (gcal_event_get_date_end (changed_event), gcal_event_get_date_start (changed_event));

  /*
   * Set the event's start and end dates. Since the event may have a
   * NULL end date, so we have to check it here
   */
  gcal_event_set_all_day (changed_event, FALSE);
  gcal_event_set_date_start (changed_event, dnd_date);

  /* Setup the new end date */
  new_end = g_date_time_add (dnd_date, timespan);
  gcal_event_set_date_end (changed_event, new_end);

  /* Commit the changes */
  gcal_manager_update_event (gcal_context_get_manager (context), changed_event, mod_type);
}

static void
on_ask_recurrence_response_cb (GcalEvent             *event,
                               GcalRecurrenceModType  mod_type,
                               gpointer               user_data)
{
  DropData *data = user_data;

  if (mod_type != GCAL_RECURRENCE_MOD_NONE)
    move_event_to_cell (data->self, data->event, data->drop_cell, mod_type);

  g_clear_object (&data->event);
  g_clear_pointer (&data, g_free);
}

static gboolean
on_drop_target_drop_cb (GtkDropTarget *drop_target,
                        const GValue  *value,
                        gdouble        x,
                        gdouble        y,
                        GcalWeekGrid  *self)
{
  GcalEventWidget *event_widget;
  GcalEvent *event;
  gint cell;

  GCAL_ENTRY;

  if (!G_VALUE_HOLDS (value, GCAL_TYPE_EVENT_WIDGET))
    GCAL_RETURN (FALSE);

  cell = get_dnd_cell (self, x, y);
  event_widget = g_value_get_object (value);
  event = gcal_event_widget_get_event (event_widget);

  if (gcal_event_has_recurrence (event))
    {
      DropData *data;

      data = g_new0 (DropData, 1);
      data->self = self;
      data->event = g_object_ref (event);
      data->drop_cell = cell;

      gcal_utils_ask_recurrence_modification_type (GTK_WIDGET (self),

  changed_event = gcal_event_new_from_event (event);
  week_start = gcal_date_time_get_start_of_week (self->active_date);
  dnd_date = gcal_date_time_add_floating_minutes (week_start, cell * 30);

  /*
   * Calculate the diff between the dropped cell and the event's start date,
   * so we can update the end date accordingly.
   */
  timespan = g_date_time_difference (gcal_event_get_date_end (changed_event), gcal_event_get_date_start (changed_event));

  /*
   * Set the event's start and end dates. Since the event may have a
   * NULL end date, so we have to check it here
   */
  gcal_event_set_all_day (changed_event, FALSE);
  gcal_event_set_date_start (changed_event, dnd_date);

  /* Setup the new end date */
  new_end = g_date_time_add (dnd_date, timespan);
  gcal_event_set_date_end (changed_event, new_end);

  /* Commit the changes */
  gcal_manager_update_event (gcal_context_get_manager (context), changed_event, mod_type);
}

static void
on_ask_recurrence_response_cb (GcalEvent             *event,
                               GcalRecurrenceModType  mod_type,
                               gpointer               user_data)
{
  DropData *data = user_data;

  if (mod_type != GCAL_RECURRENCE_MOD_NONE)
    move_event_to_cell (data->self, data->event, data->drop_cell, mod_type);

  g_clear_object (&data->event);
  g_clear_pointer (&data, g_free);
}

static gboolean
on_drop_target_drop_cb (GtkDropTarget *drop_target,
                        const GValue  *value,
                        gdouble        x,
                        gdouble        y,
                        GcalWeekGrid  *self)
{
  GcalEventWidget *event_widget;
  GcalEvent *event;
  gint cell;

  GCAL_ENTRY;

  if (!G_VALUE_HOLDS (value, GCAL_TYPE_EVENT_WIDGET))
    GCAL_RETURN (FALSE);

  cell = get_dnd_cell (self, x, y);
  event_widget = g_value_get_object (value);
  event = gcal_event_widget_get_event (event_widget);

  if (gcal_event_has_recurrence (event))
    {
      DropData *data;

      data = g_new0 (DropData, 1);
      data->self = self;
      data->event = g_object_ref (event);
      data->drop_cell = cell;

      gcal_utils_ask_recurrence_modification_type (GTK_WIDGET (self),
                                                   event,
                                                   FALSE,
                                                   on_ask_recurrence_response_cb,
                                                   data);
    }
  else
    {
      move_event_to_cell (self, event, cell, GCAL_RECURRENCE_MOD_THIS_ONLY);
    }

  GCAL_RETURN (TRUE);
}

static GdkDragAction
on_drop_target_enter_cb (GtkDropTarget *drop_target,
                         gdouble        x,
                         gdouble        y,
                         GcalWeekGrid  *self)
{

  GCAL_ENTRY;

  self->dnd.cell = get_dnd_cell (self, x, y);
  gtk_widget_set_visible (self->dnd.widget, TRUE);

  GCAL_RETURN (GDK_ACTION_COPY);
}

static void
on_drop_target_leave_cb (GtkDropTarget *drop_target,
                         GcalWeekGrid  *self)
{
  GCAL_ENTRY;

  self->dnd.cell = -1;
  gtk_widget_set_visible (self->dnd.widget, FALSE);

  GCAL_EXIT;
}

static GdkDragAction
on_drop_target_motion_cb (GtkDropTarget *drop_target,
                          gdouble        x,
                          gdouble        y,
                          GcalWeekGrid  *self)
{
  GCAL_ENTRY;
  GcalEventWidget *event_widget;
  const GValue *value;
  GDateTime *start, *end;

  self->dnd.cell = get_dnd_cell (self, x, y);
  value = gtk_drop_target_get_value (drop_target);
  event_widget = g_value_get_object (value);

  start = gcal_event_widget_get_date_start (event_widget);
  end = gcal_event_widget_get_date_end (event_widget);

  self->dnd.event_minutes = g_date_time_difference (end, start) / G_TIME_SPAN_MINUTE;

  gtk_widget_queue_allocate (GTK_WIDGET (self));

  GCAL_RETURN (GDK_ACTION_COPY);
}

```

Well...Now that you skipped it, the core issue is that the drop target, on the `on_drop_target_drop_cb` function, is attached to individual day cells (`GcalMonthCell`), which is actually a container that holds internal child widgets (such as event labels). According to research done with LLMs, in GTK 4 (to understand more about the context), any child widgets inside a container capture input events and create their our bounds, which causes the child widgets to capture our drag and drop action. Since the drop is registered on the parent cell and not on childs, this makes the action be missed entirely.

You can see more on the implementation using week grid, is that it uses a documentation more on par with GNOME's own implementation for drag and drop. Specifically, the target drop function uses *x* and *y* coordinates that target the mouse cursor. Also, instead of attaching a drop target to all individual cells, it is attached to one giant canvas (`GcalWeekGrid`). You can compare the difference between the two functions, but a main difference is this:

```C
// Week implementation
static gboolean
on_drop_target_drop_cb (GtkDropTarget *drop_target,
                        const GValue  *value,
                        gdouble        x,
                        gdouble        y,
                        GcalWeekGrid  *self)
{
    // ...
    // at some point, it uses the x and y coordinate to get the correct cell:
    cell = get_dnd_cell (self, x, y);
    // ...
}
```

This work on the week cell was done on an old commit that inspired us, [here](https://gitlab.gnome.org/GNOME/gnome-calendar/-/commit/ee5acbd493ab82ae4f18080dc926b969f1789f93). This `cell` variable is not used on the month implementation. As a result, instead of relying on the widget the mouse is hovering on, the week view grabs the raw x and y and calculates the landing week widget, bypassing UI hierarquies and a child event hijacking the action.

Do note all of this analysis and suggestion work was done with the help of LLMs. GNOME ecossystem, in general, is against LLM usage (with good reasons) and explicit prohibits AI commits (and also recommends us to not do any work that uses AI instead of the official docs). But since we were not familiar with the ecossystem, it would be way harder to tackle this issue as a starter (yeah, new devs can't code anymore :/), so we limited ourselves to use only as a guide and as a search engine, alongside official docs and codebase, and all of the code work and testing was done by us.

## Solving the problem

To solve this problem, we followed the Week refactor commit approach: move the `GtkDropTarget` up to the `GcalMonthView` instead of attaching it to each month cell, aligning what happens on the Week view's architecture. You can see our work [here](https://gitlab.gnome.org/GNOME/gnome-calendar/-/merge_requests/806/diffs), but our main change was this:

```C
// added a get_date_at_position that checks the exact month cell we need to drop onto
static GDateTime *
get_date_at_position (GcalMonthView *self,
                      gdouble        x,
                      gdouble        y)
{
  gint height = gtk_widget_get_height (GTK_WIDGET (self));
  gint header_height = gtk_widget_get_height (self->header);

  gdouble grid_y = MAX (0, y - header_height);
  gdouble grid_height = height - header_height;
  gdouble row_height = grid_height / N_ROWS_PER_PAGE;

  gint row_idx = floor (grid_y / row_height);
  row_idx = CLAMP (row_idx, 0, N_ROWS_PER_PAGE - 1);

  GtkWidget *row_widget = g_ptr_array_index (self->week_rows, FIRST_VISIBLE_ROW_INDEX + row_idx);
  GtkWidget *hovered_cell = gcal_month_view_row_get_cell_at_x (GCAL_MONTH_VIEW_ROW (row_widget), x);

  return g_date_time_ref (gcal_month_cell_get_date (GCAL_MONTH_CELL (hovered_cell)));
}
...

// uses that function inside this new on_drop implementation, similar to what week view uses
static gboolean
on_drop_target_drop_cb (GtkDropTarget *drop_target,
                        const GValue  *value,
                        gdouble        x,
                        gdouble        y,
                        GcalMonthView *self)
{
  GcalEventWidget *event_widget;
  GcalEvent *event;

  GCAL_ENTRY;

  if (!G_VALUE_HOLDS (value, GCAL_TYPE_EVENT_WIDGET))
    return FALSE;

  event_widget = g_value_get_object (value);
  event = gcal_event_widget_get_event (event_widget);

  g_autoptr (GDateTime) target_date = get_date_at_position (self, x, y);

  if (gcal_event_has_recurrence (event))
    {
      MonthDropData *data = g_new0 (MonthDropData, 1);
      data->self = self;
      data->event = g_object_ref (event);
      data->target_date = g_date_time_ref (target_date);

      gcal_utils_ask_recurrence_modification_type (GTK_WIDGET (self),
                                                   event,
                                                   FALSE,
                                                   on_recurrence_response_cb,
                                                   data);
    }
  else
    {
      move_event_to_date (self, event, target_date, GCAL_RECURRENCE_MOD_THIS_ONLY);
    }

  self->dnd_cell = NULL;
  gtk_widget_set_visible (self->dnd_widget, FALSE);
  gtk_widget_queue_allocate (GTK_WIDGET (self));

  GCAL_RETURN (TRUE);
}


```
To mimick `get_dnd_cell`, we implemented a `get_date_at_position` function that calculates the exact date using math. Thus, this replaces the reliance on widget hit-boxes and enables the parent container to catch the drop event, without missing the action.

You can check our MR [here](https://gitlab.gnome.org/GNOME/gnome-calendar/-/merge_requests/806/diffs), as we also needed to implement more things related to the drag and drop functionality.

## Testing 

We did not contribute with any unit tests, sadly (kinda hard to do one here). We did some refactors along the way, but eventually we got a working result. My colleague had some problems with their GNOME builder version. But, finally, we got our working result:

Link: https://kaikycintra.com/static/gnome_calendar_working.mp4

## Submitting

As with any open-source project, we had to read how do maintainers expect contributions, commit messages and MRs. We've found some useful guides and the Calendar's own `CONTRIBUTING.md` and `HACKING.md` for styles, so that we could streamline and make maintainers' lifes much easier (and also, have a properly accepted patch, which did not occurr with my Linux Patch :/ )

After opening our MR, Calendar's CI pipeline ran and popped up some warnings on the `style-check-diff`. Weird. 

We've went through what we might've missed regarding code style, and since job artifacts were not present, we needed to run the style check locally to mimick the job. Then we found the tool `.gitlab-ci/clang-format-diffy.py`, likely what this job uses. Doing a quick read on it's usage, we could issue the command:

```bash
git diff -U0 upstream/main | .gitlab-ci/clang-format-diff.py -binary "clang-format" -p1
```

enabling us to check the exact diffs our code had from the upstream (not merged). This returned some real issues that might make maintainers not merge directly (such as whitespaces, indentations, etc), but also some false positives. As an example, `HACKING.md` has a clear rule on indentation regarding pointers and variables, but `clang-format` was flagging these and returning what was right to it's perspective:

```C
// note: + lines indicate clang-format recommendation
 static GdkDragAction
 on_drop_target_enter_cb (GtkDropTarget *drop_target,
-                         gdouble        x,
-                         gdouble        y,
+                         gdouble x,
+                         gdouble y,
                          GcalMonthView *self)
```

which clearly violates rules dicated on the project's style convention.

With this, we've issued our refactored MR, along with some fixes regarding commit message style and MR message style. 

## Wrapping up 

As of today, we did not receive any feedback regarding a direct MR, but we did receive some labels on our MR by maintainer [Jeff](https://gitlab.gnome.org/jfft). I reckon that since this introduces many new things, it'll take some time for them to review properly. We also received a *thumbs up* on our message regarding our work on code style checks by [Nick](https://gitlab.gnome.org/niklaswimmer), so at least 2 people saw our MR so far.

With that in mind, thanks for reading so far! Since I use Linux for my everyday work and gaming, I plan to keep contributing to this beautiful ecossystem. 
