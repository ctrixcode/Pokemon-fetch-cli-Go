# Project TODO List

This file tracks the tasks for building an interactive terminal application for fetching Pokémon data. We recommend using the `Bubble Tea` library for this.

---

### Foundational: Implement an Interactive UI
*   **What:** Build the application as an interactive Terminal User Interface (TUI) using the `Bubble Tea` library.
*   **Why:** This will create a much more engaging and user-friendly experience.
*   **How:**
    1.  Integrate the `Bubble Tea` library.
    2.  Create a model that holds the application state.
    3.  The initial view could be a menu with options like "Fetch All Pokémon" and "Find Pokémon by Name/Number".

### Feature: Interactive Pokémon List
*   **What:** After fetching all Pokémon, display them in a scrollable list within the TUI.
*   **Why:** This provides a much better user experience than printing a long list to the console.
*   **How:** Use a `list` component (a "Bubble") to display the Pokémon. The user should be able to navigate the list with arrow keys.

### Feature: Find Single Pokémon (by Name or Number)
*   **What:** Add a view that prompts the user to enter a Pokémon's name or number, then fetches and displays that single Pokémon's details.
*   **Why:** This adds a core lookup functionality to the tool.
*   **How:**
    1.  Create a new view with a `textinput` component.
    2.  When the user submits their input, call the PokéAPI endpoint for a single Pokémon (e.g., `https://pokeapi.co/api/v2/pokemon/{name-or-id}`).
    3.  Display the detailed information for the fetched Pokémon.

### (Optional) Feature: Export to CSV/JSON
*   **What:** From within the TUI, add an option to export the currently viewed data (either the full list or a single Pokémon) to a file.
*   **Why:** This brings the data export functionality into the new interactive paradigm.
*   **How:** Add an 'export' option to the views. When triggered, use the `encoding/csv` or `encoding/json` packages to write the data to a file.
