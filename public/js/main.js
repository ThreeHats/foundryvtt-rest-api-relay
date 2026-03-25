document.addEventListener("DOMContentLoaded", function () {
  // Get menu elements
  const loggedOutMenu = document.getElementById("logged-out-menu");
  const loggedInMenu = document.getElementById("logged-in-menu");
  
  // Tab switching functionality
  function setupTabButtons(parent) {
    const tabButtons = parent.querySelectorAll(".tab-button");
    const tabContents = document.querySelectorAll(".tab-content");
    
    tabButtons.forEach((button) => {
      button.addEventListener("click", () => {
        const tabName = button.getAttribute("data-tab");
        
        // Skip if this is the logout button or dashboard (handled separately)
        if (!tabName) return;
        
        // Allow switching between logged-in tabs (dashboard, api-keys)
        // Don't switch to dashboard/api-keys from logged-OUT menu buttons
        if ((tabName === "dashboard" || tabName === "api-keys") && parent === loggedOutMenu) return;
        
        // Deactivate all tabs
        parent.querySelectorAll(".tab-button").forEach((btn) => btn.classList.remove("active"));
        tabContents.forEach((content) => content.classList.remove("active"));
        
        // Activate the selected tab
        button.classList.add("active");
        document.getElementById(tabName).classList.add("active");
      });
    });
  }
  
  // Set up tab switching for both menus
  setupTabButtons(loggedOutMenu);
  setupTabButtons(loggedInMenu);
  
  // Helper to show a specific tab programmatically
  function showTab(tabId) {
    const tabContents = document.querySelectorAll(".tab-content");
    tabContents.forEach((content) => content.classList.remove("active"));
    document.getElementById(tabId).classList.add("active");
  }

  // Check URL for reset token
  const urlParams = new URLSearchParams(window.location.search);
  const resetToken = urlParams.get("reset-token");

  if (resetToken) {
    // Validate the token before showing the reset form
    fetch(`/auth/validate-reset-token/${encodeURIComponent(resetToken)}`)
      .then(r => r.json())
      .then(data => {
        if (data.valid) {
          showTab("reset-password");
          document.getElementById("reset-password-form").dataset.token = resetToken;
        } else {
          showTab("login");
          const msg = document.getElementById("login-message");
          msg.textContent = "This password reset link is invalid or has expired. Please request a new one.";
          msg.className = "message error";
        }
        // Clean the URL
        window.history.replaceState({}, document.title, window.location.pathname);
      })
      .catch(() => {
        showTab("login");
        window.history.replaceState({}, document.title, window.location.pathname);
      });
  } else {
    // Check if user is already logged in (from localStorage)
    const userData = JSON.parse(localStorage.getItem("foundryApiUser"));
    if (userData) {
      // First show dashboard with cached data and switch to logged-in menu
      showDashboard(userData);
      switchToLoggedInMenu();

      // Then fetch fresh data
      fetchUserData(userData.apiKey);
    }
  }
  
  // Function to switch to logged-in menu
  function switchToLoggedInMenu() {
    loggedOutMenu.style.display = "none";
    loggedInMenu.style.display = "inherit";
  }
  
  // Function to switch to logged-out menu
  function switchToLoggedOutMenu() {
    loggedInMenu.style.display = "none";
    loggedOutMenu.style.display = "inherit";
  }
  
  // Function to fetch fresh user data
  async function fetchUserData(apiKey) {
    try {
      const response = await fetch("/auth/user-data", {
        method: "GET",
        headers: {
          "x-api-key": apiKey,
        },
      });

      if (response.ok) {
        const freshData = await response.json();

        // Fetch subscription status separately
        try {
          const subResponse = await fetch("/api/subscriptions/status", {
            method: "GET",
            headers: {
              "x-api-key": apiKey,
            },
          });
          
          if (subResponse.ok) {
            const subData = await subResponse.json();
            freshData.subscriptionStatus = subData.subscriptionStatus;
            freshData.subscriptionEndsAt = subData.subscriptionEndsAt;
          }
        } catch (subError) {
          console.error("Failed to fetch subscription status:", subError);
        }

        // Update localStorage with fresh data
        localStorage.setItem("foundryApiUser", JSON.stringify(freshData));

        // Update dashboard with fresh data
        updateDashboardData(freshData);
      }
    } catch (error) {
      console.error("Failed to fetch fresh user data:", error);
    }
  }

  // Function to update dashboard data
  function updateDashboardData(userData) {
    document.getElementById("user-email").textContent = userData.email;
    
    // Store the full API key in a data attribute and display masked version
    const apiKeyElement = document.getElementById("user-api-key");
    apiKeyElement.dataset.fullKey = userData.apiKey;
    apiKeyElement.dataset.masked = "true";
    apiKeyElement.textContent = maskApiKey(userData.apiKey);
    
    // Update rate limits display
    if (userData.limits) {
      const rateLimitsEl = document.querySelector(".message.warning");
      if (rateLimitsEl) {
        const isUnlimited = userData.limits.unlimitedMonthly;
        const dailyLimit = userData.limits.dailyLimit;
        const monthlyLimit = userData.limits.monthlyLimit;
        
        rateLimitsEl.innerHTML = `
          <strong>⚠️ Rate Limits:</strong> All users are limited to ${dailyLimit.toLocaleString()} requests per day. 
          ${isUnlimited ? 
            'You have unlimited monthly access with your subscription.' : 
            `Free accounts are limited to ${monthlyLimit} requests per month. Subscribe for unlimited monthly access.`
          }
        `;
      }
    }
    
    let status = userData.subscriptionStatus || 'free';
    // Update subscription UI
    updateSubscriptionUI(status);
    
    // Update request counts display
    const requestsToday = userData.requestsToday || 0;
    const requestsThisMonth = userData.requestsThisMonth || 0;
    const limits = userData.limits || {};
    
    if (status !== 'active') {
      document.getElementById("user-requests").textContent =
        `Monthly: ${requestsThisMonth} / ${limits.monthlyLimit || 100}, Daily: ${requestsToday} / ${limits.dailyLimit || 1000}`;
    } else {
      document.getElementById("user-requests").textContent = 
        `Monthly: ${requestsThisMonth} (unlimited), Daily: ${requestsToday} / ${limits.dailyLimit || 1000}`;
    }
  }

  // Handle signup form submission
  const signupForm = document.getElementById("signup-form");
  signupForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const email = document.getElementById("signup-email").value;
    const password = document.getElementById("signup-password").value;
    const messageEl = document.getElementById("signup-message");

    try {
      const response = await fetch("/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok) {
        messageEl.textContent = "Account created successfully!";
        messageEl.className = "message success";

        // Save user data, show dashboard, and switch to logged-in menu
        localStorage.setItem("foundryApiUser", JSON.stringify(data));
        showDashboard(data);
        switchToLoggedInMenu();
      } else {
        messageEl.textContent = data.error || "Failed to create account.";
        messageEl.className = "message error";
      }
    } catch (error) {
      messageEl.textContent = "An error occurred. Please try again.";
      messageEl.className = "message error";
      console.error(error);
    }
  });

  // Handle login form submission
  const loginForm = document.getElementById("login-form");
  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;
    const messageEl = document.getElementById("login-message");

    try {
      const response = await fetch("/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok) {
        messageEl.textContent = "Login successful!";
        messageEl.className = "message success";

        // Save user data, show dashboard, and switch to logged-in menu
        localStorage.setItem("foundryApiUser", JSON.stringify(data));
        showDashboard(data);
        switchToLoggedInMenu();
      } else {
        messageEl.textContent = data.error || "Invalid credentials.";
        messageEl.className = "message error";
      }
    } catch (error) {
      messageEl.textContent = "An error occurred. Please try again.";
      messageEl.className = "message error";
      console.error(error);
    }
  });

  // Handle logout
  const logoutBtn = document.getElementById("logout-btn");
  logoutBtn.addEventListener("click", () => {
    localStorage.removeItem("foundryApiUser");
    
    // Switch to logged-out menu
    switchToLoggedOutMenu();

    // Show signup tab
    const signupButton = document.querySelector('[data-tab="signup"]');
    signupButton.classList.add("active");
    document.getElementById("signup").classList.add("active");
    
    // Hide dashboard tab
    document.getElementById("dashboard").classList.remove("active");

    // Clear forms
    document.getElementById("signup-form").reset();
    document.getElementById("login-form").reset();
    document.getElementById("signup-message").textContent = "";
    document.getElementById("login-message").textContent = "";
  });

  // Helper function to mask API key
  function maskApiKey(apiKey) {
    if (!apiKey || apiKey.length < 12) return '••••••••••••';
    // Show first 8 and last 4 characters
    return apiKey.substring(0, 8) + '••••••••••••' + apiKey.substring(apiKey.length - 4);
  }

  // Toggle API key visibility
  const toggleApiKeyBtn = document.getElementById("toggle-api-key");
  toggleApiKeyBtn.addEventListener("click", () => {
    const apiKeyElement = document.getElementById("user-api-key");
    const isMasked = apiKeyElement.dataset.masked === "true";
    
    if (isMasked) {
      apiKeyElement.textContent = apiKeyElement.dataset.fullKey;
      apiKeyElement.dataset.masked = "false";
      toggleApiKeyBtn.textContent = "Hide";
    } else {
      apiKeyElement.textContent = maskApiKey(apiKeyElement.dataset.fullKey);
      apiKeyElement.dataset.masked = "true";
      toggleApiKeyBtn.textContent = "Show";
    }
  });

  // Copy API key to clipboard
  const copyApiKeyBtn = document.getElementById("copy-api-key");
  copyApiKeyBtn.addEventListener("click", () => {
    const apiKeyElement = document.getElementById("user-api-key");
    const apiKey = apiKeyElement.dataset.fullKey || apiKeyElement.textContent;
    navigator.clipboard.writeText(apiKey).then(() => {
      const originalText = copyApiKeyBtn.textContent;
      copyApiKeyBtn.textContent = "Copied!";
      setTimeout(() => {
        copyApiKeyBtn.textContent = originalText;
      }, 1500);
    });
  });

  // Regenerate API key
  const regenApiKeyBtn = document.getElementById("regen-api-key");
  regenApiKeyBtn.addEventListener("click", async () => {
    if (!confirm("Are you sure you want to regenerate your API key? This will invalidate your current key and you'll need to update any applications using it.")) {
      return;
    }

    // Get user data from localStorage for email
    const userData = JSON.parse(localStorage.getItem("foundryApiUser"));
    if (!userData || !userData.email) {
      alert("Please log in again to regenerate your API key.");
      return;
    }

    // Prompt for password
    const password = prompt("Please enter your password to confirm:");
    if (!password) {
      return;
    }

    try {
      regenApiKeyBtn.disabled = true;
      regenApiKeyBtn.textContent = "Regenerating...";

      const response = await fetch("/auth/regenerate-key", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: userData.email,
          password: password
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // Update the displayed API key (masked)
        const apiKeyElement = document.getElementById("user-api-key");
        apiKeyElement.dataset.fullKey = data.apiKey;
        apiKeyElement.dataset.masked = "true";
        apiKeyElement.textContent = maskApiKey(data.apiKey);
        document.getElementById("toggle-api-key").textContent = "Show";
        
        // Update localStorage with new API key
        userData.apiKey = data.apiKey;
        localStorage.setItem("foundryApiUser", JSON.stringify(userData));
        
        // Show success message
        regenApiKeyBtn.textContent = "Generated!";
        setTimeout(() => {
          regenApiKeyBtn.textContent = "Regenerate";
        }, 2000);
        
        alert("API key regenerated successfully! Please update any applications using the old key.");
      } else {
        alert(data.error || "Failed to regenerate API key");
      }
    } catch (error) {
      console.error("Error regenerating API key:", error);
      alert("An error occurred while regenerating the API key. Please try again.");
    } finally {
      regenApiKeyBtn.disabled = false;
      if (regenApiKeyBtn.textContent === "Regenerating...") {
        regenApiKeyBtn.textContent = "Regenerate";
      }
    }
  });

  // Function to show dashboard
  function showDashboard(userData) {
    // Hide all tabs and show dashboard
    const tabContents = document.querySelectorAll(".tab-content");
    tabContents.forEach((content) => content.classList.remove("active"));
    document.getElementById("dashboard").classList.add("active");
    
    // If we're showing the dashboard, also make sure dashboard tab is active in the logged-in menu
    const dashboardTab = loggedInMenu.querySelector('[data-tab="dashboard"]');
    if (dashboardTab) {
      loggedInMenu.querySelectorAll(".tab-button").forEach(btn => btn.classList.remove("active"));
      dashboardTab.classList.add("active");
    }

    // Populate user data using the shared function
    updateDashboardData(userData);
  }
  
  // Event handler for subscription button
  const subscribeBtn = document.getElementById("subscribe-btn");
  const manageSubscriptionBtn = document.getElementById("manage-subscription-btn");
  
  if (subscribeBtn) {
    subscribeBtn.addEventListener("click", async () => {
      try {
        console.log("Subscribe button clicked");
        const userData = JSON.parse(localStorage.getItem("foundryApiUser"));
        console.log("User data from localStorage:", userData);
        
        if (!userData || !userData.apiKey) {
          alert("Please log in first");
          return;
        }
        
        const response = await fetch("/api/subscriptions/create-checkout-session", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "x-api-key": userData.apiKey
          }
        });
        
        console.log("Response status:", response.status);
        
        const responseData = await response.json().catch(e => ({ error: "Failed to parse JSON" }));
        console.log("Response data:", responseData);
        
        if (response.ok) {
          window.location = responseData.url;
        } else {
          alert(responseData.error || "Failed to create checkout session");
        }
      } catch (error) {
        console.error("Detailed error:", error);
        alert("An error occurred. Please check console for details.");
      }
    });
  }
  
  if (manageSubscriptionBtn) {
    manageSubscriptionBtn.addEventListener("click", async () => {
      const userData = JSON.parse(localStorage.getItem("foundryApiUser"));
      if (!userData || !userData.apiKey) {
        alert("Please log in first");
        return;
      }
      
      try {
        const response = await fetch("/api/subscriptions/create-portal-session", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "x-api-key": userData.apiKey
          }
        });
        
        if (response.ok) {
          const { url } = await response.json();
          window.location = url;
        } else {
          alert("Failed to access subscription management");
        }
      } catch (error) {
        console.error("Error creating portal session:", error);
        alert("An error occurred. Please try again.");
      }
    });
  }
  
  // Export Data button handler (GDPR/CCPA)
  const exportDataBtn = document.getElementById("export-data-btn");
  if (exportDataBtn) {
    exportDataBtn.addEventListener("click", async () => {
      const userData = JSON.parse(localStorage.getItem("foundryApiUser"));
      if (!userData || !userData.apiKey) {
        alert("Please log in first");
        return;
      }
      
      try {
        const response = await fetch("/auth/export-data", {
          method: "GET",
          headers: {
            "x-api-key": userData.apiKey
          }
        });
        
        if (response.ok) {
          const data = await response.json();
          // Download as JSON file
          const blob = new Blob([JSON.stringify(data, null, 2)], { type: "application/json" });
          const url = URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = url;
          a.download = "my-foundry-api-data.json";
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          URL.revokeObjectURL(url);
        } else {
          alert("Failed to export data. Please try again.");
        }
      } catch (error) {
        console.error("Error exporting data:", error);
        alert("An error occurred. Please try again.");
      }
    });
  }
  
  // Delete Account button handler (GDPR/CCPA)
  const deleteAccountBtn = document.getElementById("delete-account-btn");
  if (deleteAccountBtn) {
    deleteAccountBtn.addEventListener("click", async () => {
      const userData = JSON.parse(localStorage.getItem("foundryApiUser"));
      if (!userData || !userData.apiKey) {
        alert("Please log in first");
        return;
      }
      
      // First confirmation
      const confirmed = confirm(
        "Are you sure you want to delete your account?\n\n" +
        "This action is PERMANENT and cannot be undone.\n" +
        "All your data will be deleted."
      );
      
      if (!confirmed) return;
      
      // Email verification
      const confirmEmail = prompt("Please enter your email address to confirm account deletion:");
      if (!confirmEmail) return;
      
      if (confirmEmail !== userData.email) {
        alert("Email address does not match. Account deletion cancelled.");
        return;
      }
      
      // Password verification
      const password = prompt("Please enter your password to confirm account deletion:");
      if (!password) return;
      
      // Final confirmation
      const finalConfirm = confirm(
        "FINAL WARNING: Your account will be permanently deleted.\n\n" +
        "Click OK to proceed with deletion."
      );
      
      if (!finalConfirm) return;
      
      try {
        const response = await fetch("/auth/account", {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            "x-api-key": userData.apiKey
          },
          body: JSON.stringify({ confirmEmail, password })
        });
        
        if (response.ok) {
          alert("Your account has been deleted. We're sorry to see you go.");
          // Clear local storage and switch to logged-out menu
          localStorage.removeItem("foundryApiUser");
          switchToLoggedOutMenu();
          
          // Show login tab
          const loginButton = document.querySelector('[data-tab="login"]');
          const tabContents = document.querySelectorAll(".tab-content");
          document.querySelectorAll(".tab-button").forEach((btn) => btn.classList.remove("active"));
          tabContents.forEach((content) => content.classList.remove("active"));
          loginButton.classList.add("active");
          document.getElementById("login").classList.add("active");
        } else {
          const error = await response.json();
          alert(error.error || "Failed to delete account. Please try again.");
        }
      } catch (error) {
        console.error("Error deleting account:", error);
        alert("An error occurred. Please try again.");
      }
    });
  }
  
  // Toggle change password form
  const toggleChangePasswordBtn = document.getElementById("toggle-change-password");
  const changePasswordContainer = document.getElementById("change-password-form-container");
  if (toggleChangePasswordBtn) {
    toggleChangePasswordBtn.addEventListener("click", () => {
      const isHidden = changePasswordContainer.style.display === "none";
      changePasswordContainer.style.display = isHidden ? "block" : "none";
      toggleChangePasswordBtn.innerHTML = isHidden ? "Change Password &#9660;" : "Change Password &#9654;";
    });
  }

  // Change password form
  const changePasswordForm = document.getElementById("change-password-form");
  if (changePasswordForm) {
    changePasswordForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const currentPassword = document.getElementById("current-password").value;
      const newPassword = document.getElementById("new-password").value;
      const confirmNewPassword = document.getElementById("confirm-new-password").value;
      const messageEl = document.getElementById("change-password-message");

      if (newPassword !== confirmNewPassword) {
        messageEl.textContent = "New passwords do not match.";
        messageEl.className = "message error";
        return;
      }

      const validationError = validatePassword(newPassword);
      if (validationError) {
        messageEl.textContent = validationError;
        messageEl.className = "message error";
        return;
      }

      const userData = JSON.parse(localStorage.getItem("foundryApiUser"));
      if (!userData || !userData.apiKey) {
        messageEl.textContent = "Please log in first.";
        messageEl.className = "message error";
        return;
      }

      try {
        const response = await fetch("/auth/change-password", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "x-api-key": userData.apiKey,
          },
          body: JSON.stringify({ currentPassword, newPassword }),
        });

        const data = await response.json();

        if (response.ok) {
          messageEl.textContent = "Password changed successfully!";
          messageEl.className = "message success";
          changePasswordForm.reset();
        } else {
          messageEl.textContent = data.error || "Failed to change password.";
          messageEl.className = "message error";
        }
      } catch (error) {
        messageEl.textContent = "An error occurred. Please try again.";
        messageEl.className = "message error";
        console.error(error);
      }
    });
  }

  // Client-side password validation
  function validatePassword(password) {
    if (password.length < 8) return "Password must be at least 8 characters long";
    if (!/[A-Z]/.test(password)) return "Password must contain at least one uppercase letter";
    if (!/[a-z]/.test(password)) return "Password must contain at least one lowercase letter";
    if (!/[0-9]/.test(password)) return "Password must contain at least one number";
    return null;
  }

  // Forgot password link
  const forgotPasswordLink = document.getElementById("forgot-password-link");
  if (forgotPasswordLink) {
    forgotPasswordLink.addEventListener("click", (e) => {
      e.preventDefault();
      showTab("forgot-password");
    });
  }

  // Back to login link
  const backToLoginLink = document.getElementById("back-to-login-link");
  if (backToLoginLink) {
    backToLoginLink.addEventListener("click", (e) => {
      e.preventDefault();
      showTab("login");
    });
  }

  // Forgot password form
  const forgotPasswordForm = document.getElementById("forgot-password-form");
  if (forgotPasswordForm) {
    forgotPasswordForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const email = document.getElementById("forgot-email").value;
      const messageEl = document.getElementById("forgot-password-message");

      try {
        const response = await fetch("/auth/forgot-password", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email }),
        });

        const data = await response.json();
        messageEl.textContent = data.message || "If an account with that email exists, a password reset link has been sent.";
        messageEl.className = "message success";
      } catch (error) {
        messageEl.textContent = "An error occurred. Please try again.";
        messageEl.className = "message error";
      }
    });
  }

  // Reset password form
  const resetPasswordForm = document.getElementById("reset-password-form");
  if (resetPasswordForm) {
    resetPasswordForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const password = document.getElementById("reset-new-password").value;
      const confirmPassword = document.getElementById("reset-confirm-password").value;
      const messageEl = document.getElementById("reset-password-message");
      const token = resetPasswordForm.dataset.token;

      if (password !== confirmPassword) {
        messageEl.textContent = "Passwords do not match.";
        messageEl.className = "message error";
        return;
      }

      const validationError = validatePassword(password);
      if (validationError) {
        messageEl.textContent = validationError;
        messageEl.className = "message error";
        return;
      }

      try {
        const response = await fetch("/auth/reset-password", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ token, password }),
        });

        const data = await response.json();

        if (response.ok) {
          messageEl.textContent = "Password reset successfully! Redirecting to login...";
          messageEl.className = "message success";
          setTimeout(() => {
            showTab("login");
            const loginMsg = document.getElementById("login-message");
            loginMsg.textContent = "Password reset successfully. Please sign in with your new password.";
            loginMsg.className = "message success";
          }, 2000);
        } else {
          messageEl.textContent = data.error || "Failed to reset password.";
          messageEl.className = "message error";
        }
      } catch (error) {
        messageEl.textContent = "An error occurred. Please try again.";
        messageEl.className = "message error";
      }
    });
  }

  // Show subscription UI based on status
  function updateSubscriptionUI(status) {
    const statusElement = document.getElementById("user-subscription-status");
    const subscribeBtn = document.getElementById("subscribe-btn");
    const manageSubscriptionBtn = document.getElementById("manage-subscription-btn");
    
    // Update status display
    statusElement.textContent = status === 'active' 
      ? '✅ Active' 
      : status === 'past_due'
        ? '⚠️ Past Due'
        : '🔸 Free';
    
    // Show/hide subscription buttons
    if (status === 'active' || status === 'past_due') {
      subscribeBtn.style.display = 'none';
      manageSubscriptionBtn.style.display = 'inline-block';
    } else {
      subscribeBtn.style.display = 'inline-block';
      manageSubscriptionBtn.style.display = 'none';
    }
  }

  // ==================== Scoped API Key Management ====================

  // Fetch and refresh when api-keys tab is shown
  const apiKeysTabBtn = loggedInMenu.querySelector('[data-tab="api-keys"]');
  if (apiKeysTabBtn) {
    apiKeysTabBtn.addEventListener("click", () => {
      fetchScopedKeys();
      fetchConnectedClients();
      updateMasterKeyDisplay();
    });
  }

  // Toggle credentials section
  const toggleCredentials = document.getElementById("toggle-credentials");
  if (toggleCredentials) {
    toggleCredentials.addEventListener("click", () => {
      const fields = document.getElementById("credentials-fields");
      const isHidden = fields.style.display === "none";
      fields.style.display = isHidden ? "block" : "none";
      toggleCredentials.innerHTML = isHidden
        ? "Foundry Credentials (optional) &#9660;"
        : "Foundry Credentials (optional) &#9654;";
    });
  }

  function getApiKey() {
    const userData = JSON.parse(localStorage.getItem("foundryApiUser") || "{}");
    return userData.apiKey;
  }

  function updateMasterKeyDisplay() {
    const el = document.getElementById("master-key-display");
    if (!el) return;
    const key = getApiKey();
    if (key) {
      el.setAttribute("data-full-key", key);
      el.textContent = key.substring(0, 8) + "...";
      el.setAttribute("data-masked", "true");
    }
  }

  // Expose to global scope for inline onclick handlers
  window.toggleMasterKeyDisplay = function () {
    const el = document.getElementById("master-key-display");
    if (!el) return;
    const masked = el.getAttribute("data-masked") === "true";
    if (masked) {
      el.textContent = el.getAttribute("data-full-key") || "";
      el.setAttribute("data-masked", "false");
    } else {
      el.textContent = (el.getAttribute("data-full-key") || "").substring(0, 8) + "...";
      el.setAttribute("data-masked", "true");
    }
  };

  window.copyMasterKey = function () {
    const key = getApiKey();
    if (key) navigator.clipboard.writeText(key).then(() => alert("Master key copied!"));
  };

  async function fetchScopedKeys() {
    const apiKey = getApiKey();
    if (!apiKey) return;
    try {
      const resp = await fetch("/auth/api-keys", { headers: { "x-api-key": apiKey } });
      const data = await resp.json();
      renderScopedKeys(data.keys || []);
    } catch (err) {
      document.getElementById("scoped-keys-list").innerHTML = '<p style="color: #f55;">Failed to load scoped keys.</p>';
    }
  }

  function renderScopedKeys(keys) {
    const container = document.getElementById("scoped-keys-list");
    if (!keys.length) {
      container.innerHTML = '<p style="color: #888;">No scoped keys yet. Create one to get started.</p>';
      return;
    }

    let html = '<table class="keys-table"><thead><tr>' +
      '<th>Name</th><th>Key</th><th>Status</th><th>Scopes</th><th>Daily</th><th>Expires</th><th>Creds</th><th>Actions</th>' +
      '</tr></thead><tbody>';

    for (const k of keys) {
      const status = !k.enabled ? 'disabled' : k.isExpired ? 'expired' : 'active';
      const statusClass = status === 'active' ? 'badge-active' : status === 'disabled' ? 'badge-disabled' : 'badge-expired';
      const scopes = [];
      if (k.scopedClientId) scopes.push('Client: ' + k.scopedClientId.substring(0, 12) + '...');
      if (k.scopedUserId) scopes.push('User: ' + k.scopedUserId);
      const scopeDisplay = scopes.length ? scopes.join(', ') : '<span style="color:#888">Unrestricted</span>';
      const dailyDisplay = k.dailyLimit ? `${k.requestsToday}/${k.dailyLimit}` : '<span style="color:#888">Unlimited</span>';
      const expiresDisplay = k.expiresAt ? new Date(k.expiresAt).toLocaleDateString() : '<span style="color:#888">Never</span>';
      const credsDisplay = k.hasFoundryCredentials ? '&#10003;' : '';

      html += `<tr>
        <td>${escapeHtml(k.name)}</td>
        <td><code>${escapeHtml(k.key)}</code></td>
        <td><span class="status-badge ${statusClass}">${status}</span></td>
        <td>${scopeDisplay}</td>
        <td>${dailyDisplay}</td>
        <td>${expiresDisplay}</td>
        <td style="text-align:center">${credsDisplay}</td>
        <td class="actions-cell">
          <button class="btn-small" onclick="editScopedKey(${k.id})">Edit</button>
          <button class="btn-small" onclick="toggleScopedKey(${k.id}, ${!k.enabled})">${k.enabled ? 'Disable' : 'Enable'}</button>
          <button class="btn-small btn-danger" onclick="deleteScopedKey(${k.id}, '${escapeHtml(k.name)}')">Delete</button>
        </td>
      </tr>`;
    }

    html += '</tbody></table>';
    container.innerHTML = html;
  }

  function escapeHtml(str) {
    if (!str) return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
  }

  async function fetchConnectedClients() {
    const apiKey = getApiKey();
    if (!apiKey) return;
    try {
      const resp = await fetch("/clients", { headers: { "x-api-key": apiKey } });
      const data = await resp.json();
      const select = document.getElementById("key-scoped-client");
      if (!select) return;
      // Keep the first option (unrestricted)
      select.innerHTML = '<option value="">Unrestricted (any client)</option>';
      const clients = data.clients || data || [];
      if (Array.isArray(clients)) {
        for (const c of clients) {
          const id = c.id || c.clientId || c;
          const label = c.customName ? `${c.customName} (${id})` : id;
          select.innerHTML += `<option value="${escapeHtml(String(id))}">${escapeHtml(String(label))}</option>`;
        }
      }
    } catch (err) {
      console.error("Failed to fetch connected clients:", err);
    }
  }

  window.refreshClientDropdown = async function () {
    const btn = document.querySelector('label[for="key-scoped-client"] .btn-refresh');
    if (btn) { btn.classList.add('spinning'); btn.disabled = true; }
    await fetchConnectedClients();
    if (btn) { setTimeout(() => { btn.classList.remove('spinning'); btn.disabled = false; }, 400); }
  };

  window.refreshPlayerDropdown = async function () {
    const btn = document.querySelector('label[for="key-scoped-user"] .btn-refresh');
    if (btn) { btn.classList.add('spinning'); btn.disabled = true; }
    await fetchPlayersForDropdown();
    if (btn) { setTimeout(() => { btn.classList.remove('spinning'); btn.disabled = false; }, 400); }
  };

  // When client dropdown changes, auto-fetch players for that client
  const clientSelect = document.getElementById("key-scoped-client");
  if (clientSelect) {
    clientSelect.addEventListener("change", () => {
      fetchPlayersForDropdown();
    });
  }

  async function fetchPlayersForDropdown() {
    const apiKey = getApiKey();
    if (!apiKey) return;
    const select = document.getElementById("key-scoped-user");
    if (!select) return;

    const clientId = document.getElementById("key-scoped-client").value;
    const prevValue = select.value;
    select.innerHTML = '<option value="">Unrestricted (any user)</option>';

    if (!clientId) {
      // No client selected — can't fetch players
      return;
    }

    try {
      const resp = await fetch(`/players?clientId=${encodeURIComponent(clientId)}`, {
        headers: { "x-api-key": apiKey }
      });
      const data = await resp.json();

      // The /players response has { users: [{ id, name, role, isGM, active, ... }] }
      const players = data.users || data.data || data.players || [];
      if (Array.isArray(players)) {
        const roleNames = { 0: 'None', 1: 'Player', 2: 'Trusted', 3: 'Asst GM', 4: 'GM' };
        for (const p of players) {
          const id = p.id || p.userId || p._id;
          const name = p.name || p.username || id;
          const role = roleNames[p.role] || `role ${p.role}`;
          const active = p.active ? ' - online' : '';
          const opt = document.createElement('option');
          opt.value = name; // userId param accepts name or ID
          opt.textContent = `${name} (${role}${active})`;
          select.appendChild(opt);
        }
      }

      // Restore previous selection if still valid
      if (prevValue) select.value = prevValue;
    } catch (err) {
      console.error("Failed to fetch players:", err);
    }
  }

  window.showCreateKeyForm = function () {
    document.getElementById("scoped-key-form").style.display = "block";
    document.getElementById("scoped-key-form-title").textContent = "Create Scoped Key";
    document.getElementById("save-scoped-key-btn").textContent = "Create Key";
    document.getElementById("save-scoped-key-btn").setAttribute("onclick", "saveScopedKey()");
    document.getElementById("edit-key-id").value = "";
    document.getElementById("key-name").value = "";
    document.getElementById("key-scoped-client").value = "";
    document.getElementById("key-scoped-user").value = "";
    document.getElementById("key-daily-limit").value = "";
    document.getElementById("key-expires").value = "";
    document.getElementById("key-foundry-url").value = "";
    document.getElementById("key-foundry-username").value = "";
    document.getElementById("key-foundry-password").value = "";
    document.getElementById("key-created-banner") && (document.getElementById("key-created-banner").style.display = "none");
    document.getElementById("scoped-key-form-message").textContent = "";
    updateScopedUserWarning();
    fetchConnectedClients();
  };

  function updateScopedUserWarning() {
    const warning = document.getElementById("scoped-user-warning");
    if (!warning) return;
    const userSelect = document.getElementById("key-scoped-user");
    warning.style.display = (!userSelect || !userSelect.value) ? "block" : "none";
  }

  // Listen for changes on the scoped user dropdown
  const scopedUserSelect = document.getElementById("key-scoped-user");
  if (scopedUserSelect) {
    scopedUserSelect.addEventListener("change", updateScopedUserWarning);
  }

  window.hideCreateKeyForm = function () {
    document.getElementById("scoped-key-form").style.display = "none";
  };

  window.saveScopedKey = async function () {
    const apiKey = getApiKey();
    if (!apiKey) return;
    const editId = document.getElementById("edit-key-id").value;
    const body = {
      name: document.getElementById("key-name").value,
      scopedClientId: document.getElementById("key-scoped-client").value || null,
      scopedUserId: document.getElementById("key-scoped-user").value || null,
      dailyLimit: document.getElementById("key-daily-limit").value || null,
      expiresAt: document.getElementById("key-expires").value || null,
      foundryUrl: document.getElementById("key-foundry-url").value || null,
      foundryUsername: document.getElementById("key-foundry-username").value || null,
    };
    const pwd = document.getElementById("key-foundry-password").value;
    if (pwd) body.foundryPassword = pwd;

    const msgEl = document.getElementById("scoped-key-form-message");

    if (!body.name) {
      msgEl.textContent = "Name is required.";
      msgEl.className = "message error";
      return;
    }

    try {
      let resp;
      if (editId) {
        resp = await fetch(`/auth/api-keys/${editId}`, {
          method: "PATCH",
          headers: { "Content-Type": "application/json", "x-api-key": apiKey },
          body: JSON.stringify(body)
        });
      } else {
        resp = await fetch("/auth/api-keys", {
          method: "POST",
          headers: { "Content-Type": "application/json", "x-api-key": apiKey },
          body: JSON.stringify(body)
        });
      }

      const data = await resp.json();

      if (!resp.ok) {
        msgEl.textContent = data.error || "Failed to save key.";
        msgEl.className = "message error";
        return;
      }

      if (!editId && data.key) {
        // Show the created key
        document.getElementById("key-created-banner").style.display = "block";
        document.getElementById("created-key-value").textContent = data.key;
      }

      hideCreateKeyForm();
      fetchScopedKeys();
    } catch (err) {
      msgEl.textContent = "Network error.";
      msgEl.className = "message error";
    }
  };

  window.copyCreatedKey = function () {
    const val = document.getElementById("created-key-value").textContent;
    navigator.clipboard.writeText(val).then(() => alert("Key copied!"));
  };

  window.editScopedKey = async function (id) {
    const apiKey = getApiKey();
    if (!apiKey) return;
    try {
      const resp = await fetch("/auth/api-keys", { headers: { "x-api-key": apiKey } });
      const data = await resp.json();
      const key = (data.keys || []).find(k => k.id === id);
      if (!key) return alert("Key not found.");

      document.getElementById("scoped-key-form").style.display = "block";
      document.getElementById("scoped-key-form-title").textContent = "Edit Scoped Key";
      document.getElementById("save-scoped-key-btn").textContent = "Update Key";
      document.getElementById("edit-key-id").value = id;
      document.getElementById("key-name").value = key.name || "";
      document.getElementById("key-daily-limit").value = key.dailyLimit || "";
      document.getElementById("key-expires").value = key.expiresAt ? new Date(key.expiresAt).toISOString().slice(0, 16) : "";
      document.getElementById("key-foundry-url").value = key.foundryUrl || "";
      document.getElementById("key-foundry-username").value = key.foundryUsername || "";
      document.getElementById("key-foundry-password").value = "";
      document.getElementById("scoped-key-form-message").textContent = "";

      // Set client dropdown, then fetch players for that client, then set user
      await fetchConnectedClients();
      document.getElementById("key-scoped-client").value = key.scopedClientId || "";
      if (key.scopedClientId) {
        await fetchPlayersForDropdown();
      }
      document.getElementById("key-scoped-user").value = key.scopedUserId || "";
    } catch (err) {
      alert("Failed to load key for editing.");
    }
  };

  window.toggleScopedKey = async function (id, enabled) {
    const apiKey = getApiKey();
    if (!apiKey) return;
    try {
      await fetch(`/auth/api-keys/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json", "x-api-key": apiKey },
        body: JSON.stringify({ enabled })
      });
      fetchScopedKeys();
    } catch (err) {
      alert("Failed to update key.");
    }
  };

  window.deleteScopedKey = async function (id, name) {
    if (!confirm(`Delete scoped key "${name}"? This cannot be undone.`)) return;
    const apiKey = getApiKey();
    if (!apiKey) return;
    try {
      await fetch(`/auth/api-keys/${id}`, {
        method: "DELETE",
        headers: { "x-api-key": apiKey }
      });
      fetchScopedKeys();
    } catch (err) {
      alert("Failed to delete key.");
    }
  };
});
