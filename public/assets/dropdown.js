document.addEventListener("DOMContentLoaded", () => {
  const dropdowns = document.querySelectorAll(".dropdown-menu");

  dropdowns.forEach((dropdown) => {
    const trigger = dropdown.querySelector('[aria-haspopup="menu"]');
    const popover = dropdown.querySelector("[data-popover]");

    if (!trigger || !popover) return;

    trigger.addEventListener("click", (e) => {
      e.stopPropagation();
      const isOpen = trigger.getAttribute("aria-expanded") === "true";

      closeAllDropdowns();

      if (!isOpen) {
        openDropdown(trigger, popover);
      }
    });

    document.addEventListener("click", (e) => {
      if (!dropdown.contains(e.target)) {
        closeDropdown(trigger, popover);
      }
    });

    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        closeDropdown(trigger, popover);
      }
    });
  });

  function openDropdown(trigger, popover) {
    trigger.setAttribute("aria-expanded", "true");
    popover.setAttribute("aria-hidden", "false");
  }

  function closeDropdown(trigger, popover) {
    trigger.setAttribute("aria-expanded", "false");
    popover.setAttribute("aria-hidden", "true");
  }

  function closeAllDropdowns() {
    dropdowns.forEach((dropdown) => {
      const trigger = dropdown.querySelector('[aria-haspopup="menu"]');
      const popover = dropdown.querySelector("[data-popover]");
      if (trigger && popover) {
        closeDropdown(trigger, popover);
      }
    });
  }
});
